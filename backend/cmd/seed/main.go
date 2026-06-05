package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"onboardly-backend/internal/auth"
	"onboardly-backend/internal/db"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	if err := db.InitDB(databaseURL); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.DB.Close()

	fmt.Println("Starting to seed database...")

	// 1. Seed Users
	fmt.Println("Seeding users...")
	users := []struct {
		Email    string
		Password string
		Role     string
	}{
		{"joao.analista@onboardly.com", "Senh@123", "Analista"},
		{"maria.analista@onboardly.com", "Senh@123", "Analista"},
	}

	for _, u := range users {
		// Ignore errors (like duplicate emails)
		auth.CreateUser(u.Email, u.Password, u.Role)
	}

	// 2. Seed Clients
	fmt.Println("Seeding clients...")
	clients := []struct {
		Name string
		CNPJ string
	}{
		{"Tech Solutions Ltda", "12.345.678/0001-99"},
		{"Globex Corporation", "98.765.432/0001-11"},
		{"Acme Indústria", "11.222.333/0001-44"},
	}

	var clientIDs []string
	for _, c := range clients {
		var id string
		err := db.DB.QueryRow(`
			INSERT INTO clients (name, cnpj) 
			VALUES ($1, $2) 
			ON CONFLICT (cnpj) DO UPDATE SET name = EXCLUDED.name 
			RETURNING id`, c.Name, c.CNPJ).Scan(&id)
		if err != nil {
			fmt.Printf("Error seeding client %s: %v\n", c.Name, err)
			continue
		}
		clientIDs = append(clientIDs, id)
	}

	// 3. Seed Projects
	fmt.Println("Seeding projects...")
	if len(clientIDs) >= 3 {
		projects := []struct {
			ClientID string
			Name     string
			Status   string
		}{
			{clientIDs[0], "Implantação ERP Módulo Financeiro", "Em andamento"},
			{clientIDs[0], "Treinamento Equipe de Vendas", "Backlog"},
			{clientIDs[1], "Migração de Dados Cloud", "Go-Live"},
			{clientIDs[2], "Setup Inicial Plataforma", "Em andamento"},
		}

		for _, p := range projects {
			_, err := db.DB.Exec(`
				INSERT INTO projects (client_id, name, status, is_active) 
				VALUES ($1, $2, $3, true)
				ON CONFLICT DO NOTHING`, p.ClientID, p.Name, p.Status)
			if err != nil {
				fmt.Printf("Error seeding project %s: %v\n", p.Name, err)
			}
		}
	}

	fmt.Println("Database seeding completed successfully!")
}
