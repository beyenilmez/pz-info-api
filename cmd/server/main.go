package main

import (
	"log"
	"net/http"
	"os"

	"github.com/beyenilmez/pz-info-api/internal/server"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env (if present)
	_ = godotenv.Load()

	// Retrieve port from environment or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Set up HTTP routes
	http.HandleFunc("/", server.APIHandler)

	// Start listening
	log.Printf("Listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
