package interaction_processor

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/cloudevents/sdk-go/v2/event"
	interaction2 "github.com/rmb938/gw2groups/discord/interaction"
)

// MessagePublishedData contains the full Pub/Sub message
// See the documentation for more details:
// https://cloud.google.com/eventarc/docs/cloudevents#pubsub
type MessagePublishedData struct {
	Message PubSubMessage
}

// PubSubMessage is the payload of a Pub/Sub event.
// See the documentation for more details:
// https://cloud.google.com/pubsub/docs/reference/rest/v1/PubsubMessage
type PubSubMessage struct {
	Data []byte `json:"data"`
}

func InteractionProcessor(ctx context.Context, e event.Event) error {
	var msg MessagePublishedData
	if err := e.DataAs(&msg); err != nil {
		return fmt.Errorf("event.DataAs: %v", err)
	}

	interaction := &discordgo.Interaction{}
	err := interaction.UnmarshalJSON(msg.Message.Data)
	if err != nil {
		return fmt.Errorf("error unmarshalling message data into interaction: %w", err)
	}

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_APP_BOT_TOKEN"))
	if err != nil {
		return fmt.Errorf("error creating session client: %w", err)
	}

	log.Printf("Handling async interaction for %s\n", interaction.Type)

	err = interaction2.AsyncInteractionRouter(ctx, session, interaction)

	if err != nil {
		return fmt.Errorf("error handling interaction %s: %w", interaction.Type, err)
	}

	return nil
}
