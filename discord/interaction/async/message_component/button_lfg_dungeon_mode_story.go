package message_component

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	_const "github.com/rmb938/gw2groups/discord/const"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"github.com/rmb938/gw2groups/playfab"
	"k8s.io/utils/pointer"
)

type ButtonLFGDungeonModeStory struct{}

func (c *ButtonLFGDungeonModeStory) Handle(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error {
	playFabClient := playFabAPI.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.ServerLoginWithCustomIDRequest{
		ServerCustomId: interaction.User.ID,
		InfoRequestParameters: &playFabAPI.PlayerCombinedInfoRequestParams{
			GetUserData: pointer.Bool(true),
		},
	})
	if err != nil {
		if err != nil {
			return fmt.Errorf("error logging in with playfab customid: %w", err)
		}
	}

	var dungeonsList []string

	queuedDungeons := make([]string, 0)

	err = json.Unmarshal([]byte(loginResponse.InfoResultPayload.UserData["lfg_dungeon"].Value), &queuedDungeons)
	if err != nil {
		return fmt.Errorf("error parsing lfg_dungeon from player: %w", err)
	}

	for _, id := range _const.DungeonIDs {
		for _, selectedOption := range queuedDungeons {
			if id == selectedOption {
				dungeonsList = append(dungeonsList, _const.DungeonsIDsToName[id])
			}
		}
	}

	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playFabAPI.ServerUpdateUserDataRequest{
		Data: map[string]string{
			"lfg_dungeon_mode": "story",
		},
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	// See https://learn.microsoft.com/en-us/gaming/playfab/features/multiplayer/matchmaking/config-examples#hostsearcher-or-role-based-requirements when doing Raids
	//  Specifically the "Games may have role requirements" part for selecting healer, support, tank, ect.. (quickness, alacrity, ect..)

	gw2Client := gw2Api.NewGW2APIClient(loginResponse.InfoResultPayload.UserData["gw2-api-key"].Value)
	gw2Account, err := gw2Client.GetAccount(ctx)
	if err != nil {
		// TODO: what is gw2 api key is invalid
		return fmt.Errorf("error getting gw2 account: %w", err)
	}

	gw2WorldId := gw2Account.World
	gw2WorldLocation := gw2WorldId
	for gw2WorldLocation >= 10 {
		gw2WorldLocation = gw2WorldLocation / 10
	}

	// TODO: some sort of keep-alive
	//  if we want pure 100% discord it'll need to be a button the user
	//  clicks every now and then
	//  but there won't be a way to tell them "their time is running out"
	//  We also need something to check if their ticket has given up automatically
	//  unsure how to do that without having a server constantly running and checking :(
	//  to do this effectively we would need to spin up something for every match making request

	// TODO: start helper method to do match making ticket
	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playFabAPI.ServerUpdateUserDataRequest{
		Data: map[string]string{
			"lfg_last_ack_time": strconv.FormatInt(time.Now().Unix(), 10),
		},
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	matchMakingTicketResponse, err := playFabClient.CreateMatchMakingTicket(ctx, loginResponse.EntityToken, &playFabAPI.CreateMatchMakingTicketRequest{
		Creator: playFabAPI.MatchmakingPlayer{
			Entity: loginResponse.EntityToken.Entity,
			Attributes: playFabAPI.MatchmakingPlayerAttributes{
				DataObject: map[string]interface{}{
					"lfg_world_location": strconv.Itoa(gw2WorldLocation),
					"lfg_player_count":   1, // used for playerCount expansion rule
				},
			},
		},
		GiveUpAfterSeconds: 3599,
		QueueName:          "dungeons",
	})
	if err != nil {
		return fmt.Errorf("error creating matchmaking ticket: %w", err)
	}
	// TODO: end helper method to do match making ticket

	_, err = session.FollowupMessageEdit(interaction, interaction.Message.ID, &discordgo.WebhookEdit{
		Content: pointer.String(fmt.Sprintf("> Category: **Dungeons**\n> Dungeon: **%s**\n> Mode: **Story** __*When in story mode you are expected to watch all cutscenes and be beginner friendly*__\n \nNow matchmaking for Dungeons", strings.Join(dungeonsList, ", "))),
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Stop Matchmaking",
						CustomID: "button_reset_lfg_selection",
						Style:    discordgo.DangerButton,
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}

	message := &playfab.MatchMakingTicketMessage{
		QueueName: "dungeons",
		TicketId:  matchMakingTicketResponse.TicketId,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling MatchMakingTicketMessage: %w", err)
	}

	result := pubsubTopicPlayfabMatchmakingTickets.Publish(ctx, &pubsub.Message{
		Data: messageBytes,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("error publishing message: %w", err)
	}

	return nil

}
