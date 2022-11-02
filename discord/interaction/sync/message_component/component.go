package message_component

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Component interface {
	Handle(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) (*discordgo.InteractionResponse, error)
}
