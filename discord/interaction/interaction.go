package interaction

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type Interaction interface {
	Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error)
}
