package interaction

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type SyncInteraction interface {
	Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error)
}

type AsyncInteraction interface {
	Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) error
}
