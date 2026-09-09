// Package main boots the API application.
// It loads the environment configuration, creates the HTTP server,
// and starts listening for incoming requests.
package main

import (
	"context"
	"log"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/server"
	"github.com/Mykelsown/skinbywura_backend.git/internal/service"
	"github.com/Mykelsown/skinbywura_backend.git/internal/store"
	"github.com/joho/godotenv"
)

// main loads environment variables, initializes the server, and starts the app.
func main() {
	// Load local development settings before reading the application config.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file content: %v", err)
	}

	cfg := config.Load()

	// Establish the database connection before starting the HTTP server.
	dbStore, err := store.Connection(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	// Keep the pool open for the lifetime of the application.
	defer dbStore.Pool.Close()

	// Confirm that the database is reachable, not just that the pool was created.
	if err := dbStore.Pool.Ping(context.Background()); err != nil {
		log.Fatalf("database connection ping failed: %v", err)
	}
	log.Println("database connection verified")

	// Build the service and HTTP server, then block while they handle incoming requests.
	productService := service.New(dbStore)
	serve := server.New(cfg, productService)
	log.Println("listening on port :" + cfg.Port)
	err = serve.Run()
	if err != nil {
		log.Fatalf("failed to run server on port %s: %v", cfg.Port, err)
	}
}
