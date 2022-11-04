package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rmb938/gw2groups/discord"
	"github.com/rmb938/gw2groups/playfab"

	// load .env file
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	chiRouter := chi.NewRouter()
	chiRouter.Use(middleware.RequestID)
	chiRouter.Use(middleware.RealIP)
	chiRouter.Use(middleware.Logger)
	chiRouter.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	chiRouter.Use(middleware.Timeout(30 * time.Second))
	chiRouter.Use(middleware.Heartbeat("/ping"))

	chiRouter.Use(middleware.AllowContentType("application/json"))

	chiRouter.Mount("/discord", discord.HTTPRouter())

	chiRouter.Mount("/playfab", playfab.HTTPRouter())

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: chiRouter,
	}
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// unexpected error. port in use?
		log.Fatalf("ListenAndServe(): %v", err)
	}
}
