package message_component

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	_const "github.com/rmb938/gw2groups/discord/const"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ButtonLFGDungeonModeStory struct{}

func (c *ButtonLFGDungeonModeStory) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error) {

	playFabClient := playFabAPI.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.LoginWithCustomIDRequest{
		CustomId: pointer.String(interaction.User.ID),
		InfoRequestParameters: &playFabAPI.PlayerCombinedInfoRequestParams{
			GetUserData: pointer.Bool(true),
		},
	})
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("error logging in with playfab customid: %w", err)
		}
	}

	// for key, _ := range loginResponse.InfoResultPayload.UserData {
	//
	// }

	var dungeonsList []string // TODO: figure out how to get this list, we probably will need an internal map

	queuedDungeons := make([]string, 0)

	err = json.Unmarshal([]byte(loginResponse.InfoResultPayload.UserData["lfg_dungeon"].Value), &queuedDungeons)
	if err != nil {
		return nil, fmt.Errorf("error parsing lfg_dungeon from player: %w", err)
	}

	for key, value := range _const.Dungeons {
		for _, selectedOption := range queuedDungeons {
			if key == selectedOption {
				dungeonsList = append(dungeonsList, value)
			}
		}
	}

	_, err = playFabClient.UpdateUserData(ctx, &playFabAPI.UpdateUserDataRequest{
		Data: map[string]string{
			"lfg_dungeon_mode": "story",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error updating playfab user data: %w", err)
	}

	time.Sleep(5 * time.Second)

	// See https://learn.microsoft.com/en-us/gaming/playfab/features/multiplayer/matchmaking/config-examples#hostsearcher-or-role-based-requirements when doing Raids
	//  Specifically the "Games may have role requirements" part for selecting healer, support, tank, ect.. (quickness, alacrity, ect..)

	// TODO: do matchmaking here
	//  use https://wiki.guildwars2.com/wiki/API:2/worlds to filter to correct world, find first digit and set lfg_world=1 or =2 in ticket props

	matchMakingTicketResponse, err := playFabClient.CreateMatchMakingTicket(ctx, &playFabAPI.CreateMatchMakingTicketRequest{
		Creator: playFabAPI.MatchmakingPlayer{
			Entity: loginResponse.EntityToken.Entity,
		},
		GiveUpAfterSeconds: 3599,
		QueueName:          "dungeons",
		CustomTags: map[string]string{
			"lfg_world_location": "1", // TODO: get world location
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating matchmaking ticket: %w", err)
	}

	fmt.Printf("Matchmaking Ticket: %#v\n", matchMakingTicketResponse)

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("> Category: **Dungeons**\n> Dungeon: **%s**\n> Mode: **Story** __*When in story mode you are expected to watch all cutscenes and be beginner friendly*__\n \nNow matchmaking for Dungeons", strings.Join(dungeonsList, ", ")),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Reset LFG Selections",
							CustomID: "button_reset_lfg_selection",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},
		},
	}, nil

}
