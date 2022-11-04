package message_component

import (
	"context"
	"fmt"
	"strings"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	gw2Api "github.com/rmb938/gw2groups/pkg/gw2/api"
	playFabAPI "github.com/rmb938/gw2groups/pkg/playfab/api"
	"k8s.io/utils/pointer"
)

type ButtonResetLFGSelection struct{}

func (c *ButtonResetLFGSelection) Handle(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error {

	playFabClient := playFabAPI.NewPlayFabClient()
	loginResponse, err := playFabClient.LoginWithCustomID(ctx, &playFabAPI.ServerLoginWithCustomIDRequest{
		ServerCustomId: interaction.User.ID,
		InfoRequestParameters: &playFabAPI.PlayerCombinedInfoRequestParams{
			GetUserData: pointer.Bool(true),
		},
	})
	if err != nil {
		return fmt.Errorf("error logging in with playfab customid: %w", err)
	}

	gw2Client := gw2Api.NewGW2APIClient(loginResponse.InfoResultPayload.UserData["gw2-api-key"].Value)
	gw2Account, err := gw2Client.GetAccount(ctx)
	if err != nil {
		// TODO: what if gw2 api key is invalid
		return fmt.Errorf("error getting gw2 account: %w", err)
	}

	var lfgKeys []string
	for key, _ := range loginResponse.InfoResultPayload.UserData {
		if strings.HasPrefix(key, "lfg_") {
			lfgKeys = append(lfgKeys, key)
		}
	}

	_, err = playFabClient.UpdateUserData(ctx, loginResponse.PlayFabId, &playFabAPI.ServerUpdateUserDataRequest{
		KeysToRemove: lfgKeys,
	})
	if err != nil {
		return fmt.Errorf("error updating playfab user data: %w", err)
	}

	err = playFabClient.CancelAllMatchmakingTicketsForPlayer(ctx, loginResponse.EntityToken, &playFabAPI.CancelAllMatchmakingTicketsForPlayerRequest{
		QueueName: "dungeons",
		Entity:    loginResponse.EntityToken.Entity,
	})
	if err != nil {
		return fmt.Errorf("error canceling all dungeon match making tickets: %w", err)
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
