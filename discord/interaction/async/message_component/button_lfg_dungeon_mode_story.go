package message_component

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	_const "github.com/rmb938/gw2groups/discord/const"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ButtonLFGDungeonModeStory struct{}

func (c *ButtonLFGDungeonModeStory) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error {

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

	for key, value := range _const.Dungeons {
		for _, selectedOption := range queuedDungeons {
			if key == selectedOption {
				dungeonsList = append(dungeonsList, value)
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

	fmt.Printf("Matchmaking Ticket: %#v\n", matchMakingTicketResponse)

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

	return nil

}
