package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	// 1. Clean existing seed data (order matters due to foreign key constraints)
	fmt.Println("Cleaning existing data...")
	_, _ = db.DB.Exec(`DELETE FROM meetings`)
	_, _ = db.DB.Exec(`DELETE FROM projects`)
	_, _ = db.DB.Exec(`DELETE FROM clients`)
	// We don't delete users to avoid invalidating current login tokens, but we make sure the seed users exist.

	// 2. Seed Users
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
		auth.CreateUser(u.Email, u.Password, u.Role)
	}

	// Fetch João's ID to use as the analyst for meetings
	var joaoID string
	err = db.DB.QueryRow("SELECT id FROM users WHERE email = $1", "joao.analista@onboardly.com").Scan(&joaoID)
	if err != nil {
		log.Fatalf("Failed to retrieve seed analyst ID: %v", err)
	}

	// 3. Define dates helper
	parseTime := func(str string) time.Time {
		t, err := time.Parse("2006-01-02 15:04:05", str)
		if err != nil {
			log.Fatalf("Failed to parse date %s: %v", str, err)
		}
		return t
	}

	// 4. Seed Clients with specific created_at dates for Cohort analysis
	fmt.Println("Seeding clients...")
	clientsData := []struct {
		Name      string
		CNPJ      string
		CreatedAt string
	}{
		// April 2026 cohort (Safra Abril)
		{"Tech Solutions Ltda", "12.345.678/0001-99", "2026-04-10 10:00:00"},
		{"Globex Corporation", "98.765.432/0001-11", "2026-04-15 11:00:00"},
		{"Acme Indústria", "11.222.333/0001-44", "2026-04-20 14:00:00"},
		// May 2026 cohort (Safra Maio)
		{"Beta Desenvolvimento", "44.555.666/0001-55", "2026-05-05 09:30:00"},
		{"Gama Consulting", "77.888.999/0001-88", "2026-05-10 15:00:00"},
		{"Delta Distribuidora", "22.333.444/0001-22", "2026-05-15 16:30:00"},
		// June 2026 cohort (Safra Junho)
		{"Sigma Agro", "88.777.666/0001-33", "2026-06-01 10:00:00"},
		{"Omega Software", "11.555.999/0001-00", "2026-06-02 11:30:00"},
		{"Zeta Logistics", "33.666.999/0001-11", "2026-06-03 14:15:00"},
	}

	clientMap := make(map[string]string) // name -> id
	for _, c := range clientsData {
		var id string
		createdAt := parseTime(c.CreatedAt)
		err := db.DB.QueryRow(`
			INSERT INTO clients (name, cnpj, created_at, updated_at) 
			VALUES ($1, $2, $3, $3) 
			RETURNING id`, c.Name, c.CNPJ, createdAt).Scan(&id)
		if err != nil {
			log.Fatalf("Error seeding client %s: %v", c.Name, err)
		}
		clientMap[c.Name] = id
	}

	// 5. Seed Projects with various statuses and activated_at values
	fmt.Println("Seeding projects...")
	projectsData := []struct {
		ClientName  string
		ProjectName string
		Status      string
		ActivatedAt *time.Time
		CreatedAt   string
	}{
		// April Cohort projects
		{"Tech Solutions Ltda", "Implantação ERP Módulo Financeiro", "Go-Live", ptTime(parseTime("2026-04-12 16:00:00")), "2026-04-10 10:00:00"},
		{"Globex Corporation", "Migração de Dados Cloud", "Em andamento", ptTime(parseTime("2026-05-20 14:00:00")), "2026-04-15 11:00:00"}, // activated after 35 days (not within 30d target)
		{"Acme Indústria", "Treinamento Operacional", "Backlog", nil, "2026-04-20 14:00:00"}, // abandoned/no meeting completed

		// May Cohort projects
		{"Beta Desenvolvimento", "Setup Infraestrutura", "Go-Live", ptTime(parseTime("2026-05-06 11:30:00")), "2026-05-05 09:30:00"}, // activated within 30 days
		{"Gama Consulting", "Integração Hub de Vendas", "Em andamento", ptTime(parseTime("2026-05-15 15:30:00")), "2026-05-10 15:00:00"}, // activated within 30 days
		{"Delta Distribuidora", "Customização Relatórios", "Backlog", nil, "2026-05-15 16:30:00"}, // abandoned/no show meeting

		// June Cohort projects
		{"Sigma Agro", "Configuração CRM", "Em andamento", ptTime(parseTime("2026-06-03 16:30:00")), "2026-06-01 10:00:00"}, // activated within 30 days
		{"Omega Software", "Portal de Assinaturas", "Backlog", nil, "2026-06-02 11:30:00"}, // scheduled/not yet completed
		{"Zeta Logistics", "Dashboard de Métricas", "Backlog", nil, "2026-06-03 14:15:00"}, // scheduled/not yet completed
	}

	projectMap := make(map[string]string) // name -> id
	for _, p := range projectsData {
		var id string
		createdAt := parseTime(p.CreatedAt)
		err := db.DB.QueryRow(`
			INSERT INTO projects (client_id, name, status, is_active, activated_at, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $6) 
			RETURNING id`, clientMap[p.ClientName], p.ProjectName, p.Status, p.Status != "Go-Live", p.ActivatedAt, createdAt).Scan(&id)
		if err != nil {
			log.Fatalf("Error seeding project %s: %v", p.ProjectName, err)
		}
		projectMap[p.ProjectName] = id
	}

	// 6. Seed Meetings representing the meeting lifecycle
	fmt.Println("Seeding meetings...")
	meetingsData := []struct {
		ProjectName string
		Title       string
		ScheduledAt string
		Status      string
		CompletedAt *time.Time
		NoShow      bool
	}{
		// Tech Solutions - completed meeting, client activated on same day (first meeting activation)
		{"Implantação ERP Módulo Financeiro", "Reunião de Alinhamento Inicial", "2026-04-12 15:00:00", "completed", ptTime(parseTime("2026-04-12 16:00:00")), false},

		// Globex - completed meeting, client activated later
		{"Migração de Dados Cloud", "Kickoff Técnico", "2026-04-18 10:00:00", "completed", ptTime(parseTime("2026-04-18 11:00:00")), false},

		// Beta Desenvolvimento - completed meeting, client activated same day (first meeting activation)
		{"Setup Infraestrutura", "Reunião de Kickoff e Setup", "2026-05-06 10:00:00", "completed", ptTime(parseTime("2026-05-06 11:30:00")), false},

		// Gama Consulting - completed meeting, client activated same day
		{"Integração Hub de Vendas", "Alinhamento e Homologação", "2026-05-15 14:00:00", "completed", ptTime(parseTime("2026-05-15 15:30:00")), false},

		// Delta Distribuidora - meeting scheduled, but was a NO-SHOW (thus project was abandoned)
		{"Customização Relatórios", "Reunião de Escopo", "2026-05-18 15:00:00", "completed", ptTime(parseTime("2026-05-18 15:30:00")), true}, // marked completed but with No-Show

		// Sigma Agro - completed meeting, activated same day
		{"Configuração CRM", "Homologação Técnica", "2026-06-03 15:00:00", "completed", ptTime(parseTime("2026-06-03 16:30:00")), false},

		// Omega Software - scheduled meeting (future)
		{"Portal de Assinaturas", "Reunião de Setup", "2026-06-12 10:00:00", "scheduled", nil, false},

		// Zeta Logistics - scheduled meeting (future)
		{"Dashboard de Métricas", "Reunião de Kickoff", "2026-06-15 14:00:00", "scheduled", nil, false},
	}

	for _, m := range meetingsData {
		scheduled := parseTime(m.ScheduledAt)
		_, err := db.DB.Exec(`
			INSERT INTO meetings (project_id, analyst_id, title, scheduled_at, status, completed_at, no_show) 
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			projectMap[m.ProjectName], joaoID, m.Title, scheduled, m.Status, m.CompletedAt, m.NoShow)
		if err != nil {
			log.Fatalf("Error seeding meeting %s: %v", m.Title, err)
		}
	}

	fmt.Println("Database seeding completed successfully!")
}

func ptTime(t time.Time) *time.Time {
	return &t
}
