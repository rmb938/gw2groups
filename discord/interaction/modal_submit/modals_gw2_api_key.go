package modal_submit

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ModalsGw2ApiKey struct {
}

func (s *ModalsGw2ApiKey) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ModalSubmitInteractionData) (*discordgo.InteractionResponse, error) {
	apiKey := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	gw2Client := gw2Api.NewGW2APIClient(apiKey)

	gw2Account, err := gw2Client.GetAccount(ctx)
	if err != nil {
		if gw2Api.IsAPIError(err) {
			return &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseUpdateMessage,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Error validating API Key. Please try again: %s", err),
					Components: []discordgo.MessageComponent{
						discordgo.ActionsRow{
							Components: []discordgo.MessageComponent{
								discordgo.Button{
									Label:    "Enter API Key",
									CustomID: "button_gw2_api_key",
									Style:    discordgo.PrimaryButton,
								},
							},
						},
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("error getting gw2 account: %w", err)
	}

	playFabClient := playFabAPI.NewPlayFabClient()
	_, err = playFabClient.LoginWithCustomID(ctx, &playFabAPI.LoginWithCustomIDRequest{
		CreateAccount: pointer.Bool(true),
		CustomId:      pointer.String(interaction.User.ID),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating playfab customid: %w", err)
	}

	_, err = playFabClient.UpdateUserData(ctx, &playFabAPI.UpdateUserDataRequest{
		Data: map[string]string{
			"gw2-api-key": apiKey,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error updating playfab user data: %w", err)
	}

	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Welcome %s! Select a category to begin LFG", gw2Account.Name),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    "select_menu_lfg_category",
							Placeholder: "",
							MaxValues:   1,
							Options: []discordgo.SelectMenuOption{
								{
									Label: "Dungeons",
									Value: "dungeons",
								},
								{
									Label: "Raids",
									Value: "raids",
								},
							},
							Disabled: false,
						},
					},
				},
			},
		},
	}, nil
}
