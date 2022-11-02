package main

import (
	"log"
	"os"

	// load .env file
	_ "github.com/joho/godotenv/autoload"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/rmb938/gw2groups/discord/functions"
)

func main() {
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
