package interaction

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
)

type SyncInteraction interface {
	Handler(ctx context.Context, session *discordgo.Session, interaction *discordgo.Interaction) (*discordgo.InteractionResponse, error)
}

type AsyncInteraction interface {
	Handler(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction) error
}
