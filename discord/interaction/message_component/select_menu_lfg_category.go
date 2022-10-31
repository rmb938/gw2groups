package message_component

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type SelectMenuLFGCategory struct{}

func (c *SelectMenuLFGCategory) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error) {
	playFabClient := playFabAPI.NewPlayFabClient()
	_, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.LoginWithCustomIDRequest{
		CustomId: pointer.String(interaction.User.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("error logging in with playfab customid: %w", err)
	}

	_, err = playFabClient.UpdateUserData(ctx, &playFabAPI.UpdateUserDataRequest{
		Data: map[string]string{
			"lfg_category": data.Values[0],
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error updating playfab user data: %w", err)
	}

	switch value := data.Values[0]; value {
	case "dungeons":
		dungeonOptions := []discordgo.SelectMenuOption{
			{
				Label: "Any",
				Value: "any",
			},
			{
				Label: "Ascalonian Catacombs",
				Value: "ascalonian_catacombs",
			},
			{
				Label: "Caudecus's Manor",
				Value: "caudecuss_manor",
			},
			{
				Label: "Twilight Arbor",
				Value: "twilight_arbor",
			},
			{
				Label: "Sorrow's Embrace",
				Value: "sorrows_embrace",
			},
			{
				Label: "Citadel of Flame",
				Value: "citadel_of_flame",
			},
			{
				Label: "Honor of the Waves",
				Value: "honor_of_the_waves",
			},
			{
				Label: "Crucible of Eternity",
				Value: "crucible_of_eternity",
			},
			{
				Label: "The Ruined City of Arah",
				Value: "the_ruined_city_of_arah",
			},
		}

		return &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(">>> Category: **Dungeons**"),
				Components: []discordgo.MessageComponent{
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
			},
		}, nil
	default:
		return nil, fmt.Errorf("select menu lfg category %s is not defined", value)
	}
}
