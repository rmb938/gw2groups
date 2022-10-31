package message_component

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ButtonResetLFGSelection struct{}

func (c *ButtonResetLFGSelection) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error) {

	playFabClient := playFabAPI.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.LoginWithCustomIDRequest{
		CustomId: pointer.String(interaction.User.ID),
		InfoRequestParameters: &playFabAPI.PlayerCombinedInfoRequestParams{
			GetUserData: pointer.Bool(true),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error creating playfab customid: %w", err)
	}

	gw2Client := gw2Api.NewGW2APIClient(loginResponse.InfoResultPayload.UserData["gw2-api-key"].Value)

	gw2Account, err := gw2Client.GetAccount(ctx)
	if err != nil {
		// TODO: what is gw2 api key is invalid
		return nil, fmt.Errorf("error getting gw2 account: %w", err)
	}

	var lfgKeys []string
	for key, _ := range loginResponse.InfoResultPayload.UserData {
		if strings.HasPrefix(key, "lfg_") {
			lfgKeys = append(lfgKeys, key)
		}
	}

	_, err = playFabClient.UpdateUserData(ctx, &playFabAPI.UpdateUserDataRequest{
		KeysToRemove: lfgKeys,
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
