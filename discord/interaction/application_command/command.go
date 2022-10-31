package application_command

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.ApplicationCommandInteractionData) (*discordgo.InteractionResponse, error)
	CanDM() bool
	CanGuild() bool
}
