package config

import (
	"os"
	"testing"
)

func TestLoadConfig_Defaults(t *testing.T) {
	// Backup original env vars
	origPort := os.Getenv("PORT")
	origDBURL := os.Getenv("DATABASE_URL")
	origJWTSecret := os.Getenv("JWT_SECRET")
	origWebhookToken := os.Getenv("WEBHOOK_TOKEN")

	defer func() {
		os.Setenv("PORT", origPort)
		os.Setenv("DATABASE_URL", origDBURL)
		os.Setenv("JWT_SECRET", origJWTSecret)
		os.Setenv("WEBHOOK_TOKEN", origWebhookToken)
	}()

	// Clear environment variables to force defaults
	os.Unsetenv("PORT")
	os.Unsetenv("DATABASE_URL")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("WEBHOOK_TOKEN")

	cfg := LoadConfig()

	if cfg.Port != "8080" {
		t.Errorf("Expected Port to be '8080', got '%s'", cfg.Port)
	}
	expectedDB := "postgresql://postgres:postgres@localhost:5432/onboardlyv2?sslmode=disable"
	if cfg.DatabaseURL != expectedDB {
		t.Errorf("Expected DatabaseURL to be '%s', got '%s'", expectedDB, cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "super-secret-onboardly-key" {
		t.Errorf("Expected JWTSecret to be 'super-secret-onboardly-key', got '%s'", cfg.JWTSecret)
	}
	if cfg.WebhookToken != "your-shared-secret-token" {
		t.Errorf("Expected WebhookToken to be 'your-shared-secret-token', got '%s'", cfg.WebhookToken)
	}
}

func TestLoadConfig_Overrides(t *testing.T) {
	// Backup original env vars
	origPort := os.Getenv("PORT")
	origDBURL := os.Getenv("DATABASE_URL")
	origJWTSecret := os.Getenv("JWT_SECRET")
	origWebhookToken := os.Getenv("WEBHOOK_TOKEN")

	defer func() {
		os.Setenv("PORT", origPort)
		os.Setenv("DATABASE_URL", origDBURL)
		os.Setenv("JWT_SECRET", origJWTSecret)
		os.Setenv("WEBHOOK_TOKEN", origWebhookToken)
	}()

	// Set custom environment variables
	os.Setenv("PORT", "9090")
	os.Setenv("DATABASE_URL", "postgresql://testuser:testpass123@localhost:5433/onboardlyv2_test?sslmode=disable")
	os.Setenv("JWT_SECRET", "custom-jwt-secret")
	os.Setenv("WEBHOOK_TOKEN", "custom-webhook-token")

	cfg := LoadConfig()

	if cfg.Port != "9090" {
		t.Errorf("Expected Port to be '9090', got '%s'", cfg.Port)
	}
	expectedDB := "postgresql://testuser:testpass123@localhost:5433/onboardlyv2_test?sslmode=disable"
	if cfg.DatabaseURL != expectedDB {
		t.Errorf("Expected DatabaseURL to be '%s', got '%s'", expectedDB, cfg.DatabaseURL)
	}
	if cfg.JWTSecret != "custom-jwt-secret" {
		t.Errorf("Expected JWTSecret to be 'custom-jwt-secret', got '%s'", cfg.JWTSecret)
	}
	if cfg.WebhookToken != "custom-webhook-token" {
		t.Errorf("Expected WebhookToken to be 'custom-webhook-token', got '%s'", cfg.WebhookToken)
	}
}
