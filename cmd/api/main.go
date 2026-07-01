package main

import (
	"fmt"
	"log"

	"github.com/Mykelsown/skinbywura_backend.git/config"
	"github.com/Mykelsown/skinbywura_backend.git/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("failed to load .env file content: %v", err)
	}

	serve := server.New(config.Load())
	fmt.Println("server is running on port "+ config.Load().Port)
	err  = serve.Run()
	if err != nil {
		log.Fatal("failed to run server on port " + config.Load().Port)
	}
}