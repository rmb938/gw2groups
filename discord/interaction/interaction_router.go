package interaction

import (
	"context"
	"fmt"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/rmb938/gw2groups/discord/interaction/async"
	"github.com/rmb938/gw2groups/discord/interaction/sync"
)

var syncInteractionHandlers = map[discordgo.InteractionType]SyncInteraction{
	discordgo.InteractionPing:               &sync.Ping{},
	discordgo.InteractionApplicationCommand: &sync.ApplicationCommand{},
	discordgo.InteractionMessageComponent:   &sync.MessageComponent{},
}

var asyncInteractionHandlers = map[discordgo.InteractionType]AsyncInteraction{
	discordgo.InteractionApplicationCommand: &async.ApplicationCommand{},
	discordgo.InteractionMessageComponent:   &async.MessageComponent{},
	discordgo.InteractionModalSubmit:        &async.ModalSubmit{},
}

func SyncInteractionRouter(ctx context.Context, session *discordgo.Session, i *discordgo.Interaction, interactionBytes []byte) (*discordgo.InteractionResponse, error) {
	pubsubClient, err := pubsub.NewClient(ctx, os.Getenv("PUBSUB_PROJECT_ID"))
	if err != nil {
		return nil, fmt.Errorf("error creating pubsub client: %w", err)
	}
	pubsubTopic := pubsubClient.Topic(os.Getenv("PUBSUB_TOPIC_ID"))

	handler, ok := syncInteractionHandlers[i.Type]
	var response *discordgo.InteractionResponse
	if ok {
		response, err = handler.Handler(ctx, session, i)
		if err != nil {
			return nil, fmt.Errorf("error handling sync interaction: %w", err)
		}
	}

	// check if handled in async
	_, ok = asyncInteractionHandlers[i.Type]
	if ok {
		// is handled async to send it off
		result := pubsubTopic.Publish(ctx, &pubsub.Message{
			Data: interactionBytes,
		})

		_, err = result.Get(ctx)
		if err != nil {
			return nil, fmt.Errorf("error publishing message: %w", err)
		}

		if response == nil {
			if i.Type == discordgo.InteractionApplicationCommand {
				return &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
				}, nil
			}

			if i.Type == discordgo.InteractionMessageComponent || i.Type == discordgo.InteractionModalSubmit {
				return &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseDeferredMessageUpdate,
				}, nil
			}

			return nil, nil
		}
	}

	return response, nil
}

func AsyncInteractionRouter(ctx context.Context, session *discordgo.Session, i *discordgo.Interaction) error {
	if handler, ok := asyncInteractionHandlers[i.Type]; ok {
		// TODO: what to do with these errors
		//   we most likely can't retry a majority of them
		//   so we probably shouldn't actually error, just log
		//   if we actually error we are going to retry until pubsub times out

		err := handler.Handler(ctx, session, i)
		if err != nil {
			log.Printf("error handling async interaction: %s: %s", i.Type, err)
		}

		return nil
	}

	// TODO: what do we do if we got an async but don't handle it (i.e http updated but pubsub hasn't yet)

	return nil
}
