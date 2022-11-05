package discord

import (
	"context"
	"encoding/hex"
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
	_const "github.com/rmb938/gw2groups/discord/const"
	interaction2 "github.com/rmb938/gw2groups/discord/interaction"
)

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

	pubsubTopicDiscordInteractions := pubsubClient.Topic(os.Getenv("PUBSUB_DISCORD_INTERACTIONS_TOPIC_ID"))
	pubsubTopicPlayfabMatchmakingTickets := pubsubClient.Topic(os.Getenv("PUBSUB_PLAYFAB_MATCHMAKING_TICKETS_TOPIC_ID"))

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_APP_BOT_TOKEN"))
	if err != nil {
		panic(fmt.Errorf("error creating session client: %w", err))
	}

	// TODO: move command overwrite to it's own cmd
	//  it'll be executed at deploy time (or via makefile before the normal run)
	_, err = session.ApplicationCommandBulkOverwrite(
		os.Getenv("DISCORD_APP_ID"),
		"",
		[]*discordgo.ApplicationCommand{
			{
				Name:        "lfg",
				Description: "Have the BOT send a DM to start a LFG session",
			},
		},
	)
	if err != nil {
		panic(fmt.Errorf("error registering session commands: %w", err))
	}

	rawPublicKey := os.Getenv("DISCORD_APP_PUBLIC_KEY")
	hexPublicKey, _ := hex.DecodeString(rawPublicKey)

	chiRouter.Post("/interactions", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()

		verified := discordgo.VerifyInteraction(r, hexPublicKey)
		if !verified {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx = _const.SetPubsubClient(ctx, pubsubClient)
		ctx = _const.SetPubsubTopic(ctx, pubsubTopicDiscordInteractions)
		ctx = _const.SetPubsubTopic(ctx, pubsubTopicPlayfabMatchmakingTickets)

		bodyRaw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		interaction := &discordgo.Interaction{}
		err = interaction.UnmarshalJSON(bodyRaw)
		if err != nil {
			log.Printf("error unmarshalling interaction: %s", err)
			http.Error(w, "error unmarshalling interaction", http.StatusBadRequest)
			return
		}

		response, err := interaction2.SyncInteractionRouter(ctx, session, interaction, bodyRaw)
		if err != nil {
			log.Printf("error handling interaction %s: %s\n", interaction.Type, err)
			http.Error(w, "error handling interaction", http.StatusInternalServerError)
			return
		}

		if response == nil {
			log.Printf("interaction has no response %s\n", interaction.Type)
			http.Error(w, "interaction has no response", http.StatusNotImplemented)
			return
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, response)
	})

	chiRouter.Post("/process-interaction", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()

		if _, exists := os.LookupEnv("PUBSUB_EMULATOR_HOST"); !exists {
			// TODO: validate Authorization header
			//  subscription should be configured to send a service account
			//  no idea how to test this locally with the emulator though
		}

		ctx = _const.SetPubsubClient(ctx, pubsubClient)
		ctx = _const.SetPubsubTopic(ctx, pubsubTopicDiscordInteractions)
		ctx = _const.SetPubsubTopic(ctx, pubsubTopicPlayfabMatchmakingTickets)

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

		interaction := &discordgo.Interaction{}
		err = interaction.UnmarshalJSON(message.Message.Data)
		if err != nil {
			log.Printf("error unmarshalling interaction: %s", err)
			http.Error(w, "error unmarshalling interaction", http.StatusBadRequest)
			return
		}

		err = interaction2.AsyncInteractionRouter(ctx, session, interaction)

		if err != nil {
			log.Printf("error handling interaction %s: %s", interaction.Type, err)
			http.Error(w, "error handling interaction", http.StatusBadRequest)
		}

		render.Status(r, http.StatusNoContent)
	})

	return chiRouter
}
