package message_component

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	_const "github.com/rmb938/gw2groups/discord/const"
	"github.com/rmb938/gw2groups/pkg/api_clients/playfab"
	"k8s.io/utils/pointer"
)

type SelectMenuLFGCategory struct{}

func (c *SelectMenuLFGCategory) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error {
	playFabClient := playfab.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playfab.ServerLoginWithCustomIDRequest{
		ServerCustomId: interaction.User.ID,
	})
	if err != nil {
		return fmt.Errorf("error logging in with playfab customid: %w", err)
	}

	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playfab.ServerUpdateUserDataRequest{
		Data: map[string]string{
			"lfg_category": data.Values[0],
		},
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	switch value := data.Values[0]; value {
	case "dungeons":
		dungeonOptions := []discordgo.SelectMenuOption{
			{
				Label: "Any",
				Value: "any",
			},
		}

		for _, id := range _const.DungeonIDs {
			dungeonOptions = append(dungeonOptions, discordgo.SelectMenuOption{
				Label: _const.DungeonsIDsToName[id],
				Value: id,
			})
		}

		_, err = session.FollowupMessageEdit(interaction, interaction.Message.ID, &discordgo.WebhookEdit{
			Content: pointer.String(fmt.Sprintf("> Category: **Dungeons**\n \nWhich Dungeon(s) would you like to group for?")),
			Components: &[]discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "select_menu_lfg_dungeon",
							Placeholder: "",
							MinValues:   pointer.Int(1),
							MaxValues:   len(dungeonOptions),
							Options:     dungeonOptions,
							Disabled:    false,
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
	default:
		return fmt.Errorf("select menu lfg category %s is not defined", value)
	}
}
