package config

import (
	"os"
)

// Config holds all configuration parameters for the application.
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	WebhookToken string
}

// LoadConfig reads configuration parameters from environment variables with safe defaults.
func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	dbURL := getEnv("DATABASE_URL", "postgresql://postgres:postgres@localhost:5432/onboardlyv2?sslmode=disable")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-onboardly-key")
	webhookToken := getEnv("WEBHOOK_TOKEN", "your-shared-secret-token")

	return &Config{
		Port:         port,
		DatabaseURL:  dbURL,
		JWTSecret:    jwtSecret,
		WebhookToken: webhookToken,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
