package modal_submit

import (
	"context"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ModalsGw2ApiKey struct {
}

func (s *ModalsGw2ApiKey) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ModalSubmitInteractionData) error {
	apiKey := data.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value

	gw2Client := gw2Api.NewGW2APIClient(apiKey)

	gw2Account, err := gw2Client.GetAccount(ctx)
	if err != nil {
		if gw2Api.IsAPIError(err) {
			_, err := session.FollowupMessageEdit(interaction, interaction.Message.ID, &discordgo.WebhookEdit{
				Content: pointer.String(fmt.Sprintf("Error validating API Key. Please try again: %s", err)),
				Components: &[]discordgo.MessageComponent{
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
			})
			return err
		}
		return fmt.Errorf("error getting gw2 account: %w", err)
	}

	playFabClient := playFabAPI.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.ServerLoginWithCustomIDRequest{
		CreateAccount:  pointer.Bool(true),
		ServerCustomId: interaction.User.ID,
		InfoRequestParameters: &playFabAPI.PlayerCombinedInfoRequestParams{
			GetUserData: pointer.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("error creating playfab customid: %w", err)
	}

	var lfgKeys []string
	for key, _ := range loginResponse.InfoResultPayload.UserData {
		if strings.HasPrefix(key, "lfg_") {
			lfgKeys = append(lfgKeys, key)
		}
	}

	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playFabAPI.ServerUpdateUserDataRequest{
		Data: map[string]string{
			"gw2-api-key": apiKey,
		},
		KeysToRemove: lfgKeys,
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	_, err = session.FollowupMessageEdit(interaction, interaction.Message.ID, &discordgo.WebhookEdit{
		Content: pointer.String(fmt.Sprintf("Welcome %s! Select a category to begin LFG", gw2Account.Name)),
		Components: &[]discordgo.MessageComponent{
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
								Label: "Strike Missions",
								Value: "strike_missions",
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
	})
	if err != nil {
		return err
	}

	return nil

}
