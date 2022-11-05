package playfab

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rmb938/gw2groups/pkg/api_clients/playfab"
)

type MatchMakingTicketMessage struct {
	QueueName string `json:"QueueName"`
	TicketId  string `json:"TicketId"`
}

type PubSubHTTPMessage struct {
	Message struct {
		Attributes  map[string]string `json:"attributes"`
		Data        []byte            `json:"data,omitempty"`
		MessageID   string            `json:"message_id"`
		PublishTime time.Time         `json:"publish_time"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func HTTPRouter() *chi.Mux {
	chiRouter := chi.NewRouter()

	pubsubClient, err := pubsub.NewClient(context.TODO(), os.Getenv("PUBSUB_PROJECT_ID"))
	if err != nil {
		panic(fmt.Errorf("error creating pubsub client: %w", err))
	}

	topic := pubsubClient.Topic(os.Getenv("PUBSUB_PLAYFAB_MATCHMAKING_TICKETS_TOPIC_ID"))

	_, err = discordgo.New("Bot " + os.Getenv("DISCORD_APP_BOT_TOKEN"))
	if err != nil {
		panic(fmt.Errorf("error creating session client: %w", err))
	}

	chiRouter.Post("/matchmaking-ticket-poll", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()

		if _, exists := os.LookupEnv("PUBSUB_EMULATOR_HOST"); !exists {
			// TODO: validate Authorization header
			//  subscription should be configured to send a service account
			//  no idea how to test this locally with the emulator though
		}

		bodyRaw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		message := &PubSubHTTPMessage{}
		err = json.Unmarshal(bodyRaw, message)
		if err != nil {
			log.Printf("error unmarshalling body: %s", err)
			http.Error(w, "error unmarshalling body", http.StatusBadRequest)
			return
		}

		matchMakingTicketMessage := &MatchMakingTicketMessage{}
		err = json.Unmarshal(message.Message.Data, matchMakingTicketMessage)
		if err != nil {
			log.Printf("error unmarshalling matchmaking ticket message: %s", err)
			http.Error(w, "error unmarshalling matchmaking ticket message", http.StatusBadRequest)
			return
		}

		playFabClient := playfab.NewPlayFabClient()
		titleEntityToken, err := playFabClient.GetTitleEntityToken(ctx)
		if err != nil {
			log.Printf("error getting title entity token: %s", err)
			http.Error(w, "error getting title entity token", http.StatusBadRequest)
			return
		}

		matchMakingTicketResponse, err := playFabClient.GetMatchmakingTicket(ctx, *titleEntityToken, &playfab.GetMatchMakingTicketRequest{
			EscapeObject: false,
			QueueName:    matchMakingTicketMessage.QueueName,
			TicketId:     matchMakingTicketMessage.TicketId,
		})
		if err != nil {
			log.Printf("error getting match making ticket: %s", err)
			http.Error(w, "error getting match making ticket", http.StatusBadRequest)
			return
		}

		if matchMakingTicketResponse.Status == "Canceled" {
			render.Status(r, http.StatusNoContent)
			return
		}

		if matchMakingTicketResponse.Status == "Matched" {
			matchResponse, err := playFabClient.GetMatch(ctx, *titleEntityToken, &playfab.GetMatchRequest{
				EscapeObject:           false,
				MatchId:                *matchMakingTicketResponse.MatchId,
				QueueName:              matchMakingTicketMessage.QueueName,
				ReturnMemberAttributes: false,
				CustomTags:             nil,
			})
			if err != nil {
				log.Printf("error getting match: %s", err)
				http.Error(w, "error getting match", http.StatusBadRequest)
				return
			}

			log.Printf("%#v\n", matchResponse)

			// TODO: send message about match info

			render.Status(r, http.StatusNoContent)
			return
		}

		select {
		case <-ctx.Done():
			break
		case <-time.After(10 * time.Second):
			result := topic.Publish(ctx, &pubsub.Message{
				Data: message.Message.Data,
			})

			_, err = result.Get(ctx)
			if err != nil {
				log.Printf("error publishing message: %s", err)
				http.Error(w, "error publishing message", http.StatusBadRequest)
				return
			}
			break
		}

		render.Status(r, http.StatusNoContent)
	})

	return chiRouter
}
