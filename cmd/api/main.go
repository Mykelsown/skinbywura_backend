// Package main boots the API application.
// It loads the environment configuration, creates the HTTP server,
// and starts listening for incoming requests.
package main

import (
	"log"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/server"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/joho/godotenv"
)

// main loads environment variables, initializes the server, and starts the app.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file content: %v", err)
	}

	cfg := config.Load()
	dbStore, err := store.Connection(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbStore.Pool.Close()

	serve := server.New(cfg)
	err = serve.Run()
	if err != nil {
		log.Fatalf("failed to run server on port %s: %v", cfg.Port, err)
	}
}
