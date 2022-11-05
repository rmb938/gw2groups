package message_component

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	_const "github.com/rmb938/gw2groups/discord/const"
	"github.com/rmb938/gw2groups/pkg/api_clients/playfab"
	"k8s.io/utils/pointer"
)

type SelectMenuLFGDungeon struct{}

func (c *SelectMenuLFGDungeon) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error {
	playFabClient := playfab.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playfab.ServerLoginWithCustomIDRequest{
		ServerCustomId: interaction.User.ID,
	})
	if err != nil {
		return fmt.Errorf("error logging in with playfab customid: %w", err)
	}

	var dungeonsList []string

	hasAny := false
	for _, value := range data.Values {
		if value == "any" {
			hasAny = true
			break
		}
	}

	if hasAny {
		data.Values = []string{}
		for _, id := range _const.DungeonIDs {
			data.Values = append(data.Values, id)
		}
	}

	for _, id := range _const.DungeonIDs {
		for _, selectedOption := range data.Values {
			if id == selectedOption {
				dungeonsList = append(dungeonsList, _const.DungeonsIDsToName[id])
			}
		}
	}

	dataValues, err := json.Marshal(data.Values)
	if err != nil {
		return fmt.Errorf("error converting dungeon values to string: %w", err)
	}

	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playfab.ServerUpdateUserDataRequest{
		Data: map[string]string{
			"lfg_dungeon": string(dataValues),
		},
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	// TODO: ask if story or exploration mode
	//  if exploration as for path (any, path 1, path 2, path 3, save data, start matchmaking
	//  if story, save data, start matchmaking

	_, err = session.FollowupMessageEdit(interaction, interaction.Message.ID, &discordgo.WebhookEdit{
		Content: pointer.String(fmt.Sprintf("> Category: **Dungeons**\n> Dungeon: **%s**\n \nSelect the dungeon mode\n \n__*When in story mode you are expected to watch all cutscenes and be beginner friendly*__", strings.Join(dungeonsList, ", "))),
		Components: &[]discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Story Mode",
						CustomID: "button_lfg_dungeon_mode_story",
						Style:    discordgo.PrimaryButton,
					},
					discordgo.Button{
						Label:    "Exploration Mode",
						CustomID: "button_lfg_dungeon_mode_exploration",
						Style:    discordgo.PrimaryButton,
					},
				},
			},
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
	})
	if err != nil {
		return err
	}

	return nil
}
