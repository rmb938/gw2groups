package message_component

import (
	"context"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
)

type Component interface {
	Handle(ctx context.Context, session *discordgo.Session, pubsubTopicPlayfabMatchmakingTickets *pubsub.Topic, interaction *discordgo.Interaction, data discordgo.MessageComponentInteractionData) error
}
