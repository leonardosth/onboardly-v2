package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"onboardly-backend/internal/api"
	"onboardly-backend/internal/config"
	"onboardly-backend/internal/db"
)

func main() {
	// Load environment variables from .env file if present
	if err := godotenv.Load(); err != nil {
		log.Println("Note: No .env file found. Using default environment variables or system configuration.")
	}

	// Load environment configuration parameters
	cfg := config.LoadConfig()

	// Initialize the PostgreSQL connection pool
	log.Printf("Connecting to PostgreSQL at %s...", cfg.DatabaseURL)
	if err := db.InitDB(cfg.DatabaseURL); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.DB.Close()
	log.Println("Database connection pool established successfully.")

	// Initialize HTTP router
	router := api.NewRouter(cfg)

	// Start the API server
	log.Printf("Starting HTTP server on port %s...", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, router); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
