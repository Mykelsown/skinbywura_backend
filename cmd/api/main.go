// Package main boots the API application.
// It loads the environment configuration, creates the HTTP server,
// and starts listening for incoming requests.
package main

import (
	"fmt"
	"log"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/server"
	"github.com/joho/godotenv"
)

// main loads environment variables, initializes the server, and starts the app.
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file content: %v", err)
	}

	serve := server.New(config.Load())
	err = serve.Run()
	fmt.Println("we")
	if err != nil {
		log.Fatal("failed to run server on port " + config.Load().Port)
	}
	fmt.Println("server is running on port " + config.Load().Port)
}
