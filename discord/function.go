package discord

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

var chiRouter *chi.Mux

func init() {
	chiRouter = chi.NewRouter()
	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.RealIP)
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	chiRouter.Use(middleware.Timeout(10 * time.Second))
	chiRouter.Use(middleware.Heartbeat("/ping"))

	chiRouter.Use(middleware.AllowContentType("application/json"))

	chiRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawPublicKey := os.Getenv("DISCORD_APP_PUBLIC_KEY")
			hexPublicKey, _ := hex.DecodeString(rawPublicKey)

			verified := discordgo.VerifyInteraction(r, hexPublicKey)
			if !verified {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

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

	chiRouter.Post("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()
		bodyRaw, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "error reading body", http.StatusInternalServerError)
			return
		}

		interaction := &discordgo.Interaction{}
		err = interaction.UnmarshalJSON(bodyRaw)
		if err != nil {

			http.Error(w, "error unmarshalling body", http.StatusBadRequest)
			return
		}

		response, err := InteractionRouter(ctx, session, interaction)
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

	functions.HTTP("discord", discordInteraction)
}

func discordInteraction(w http.ResponseWriter, r *http.Request) {
	chiRouter.ServeHTTP(w, r)
}
