package interaction

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Ping struct{}

func (i *Ping) Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error) {
	return &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}, nil
}
