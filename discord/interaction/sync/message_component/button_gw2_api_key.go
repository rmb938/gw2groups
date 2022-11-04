package message_component

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type ButtonGw2ApiKey struct{}

func (c *ButtonGw2ApiKey) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error) {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "modals_gw2_api_key",
			Title:    "Enter your GW2 API Key",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:  "text_inputs_gw2_api_key",
							Label:     "GW2 API Key",
							Style:     discordgo.TextInputShort,
							Required:  true,
							MinLength: 72,
							MaxLength: 72,
						},
					},
				},
			},
		},
	}, nil
}
