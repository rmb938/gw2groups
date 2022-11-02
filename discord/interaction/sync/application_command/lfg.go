package application_command

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type LFG struct{}

func (c *LFG) Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error) {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "I'm sending you a DM.",
		},
	}, nil
}

func (c *LFG) CanDM() bool {
	return false
}

func (c *LFG) CanGuild() bool {
	return true
}
