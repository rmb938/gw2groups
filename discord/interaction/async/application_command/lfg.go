package application_command

import (
	"context"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type LFG struct{}

func (c *LFG) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) error {
	userDM, err := session.UserChannelCreate(interaction.Member.User.ID)
	if err != nil {
		return fmt.Errorf("error creating user DM channel: %w", err)
	}

	// TODO: lookup playfab player using User.ID
	//  if player has title player data gw2-api-key go right into the LFG
	//  otherwise show first time buttons (asking for api key)

	_, err = session.ChannelMessageSendComplex(userDM.ID, &discordgo.MessageSend{
		Content: "Hey there! It looks like this is your first time using LFG. Please enter your GW2 API Key", // TODO: formalize this
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
	})
	if err != nil {
		return fmt.Errorf("error sending DM message: %w", err)
	}

	return nil
}

func (c *LFG) CanDM() bool {
	return false
}

func (c *LFG) CanGuild() bool {
	return true
}
