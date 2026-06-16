package project

import "time"

// Project represents the project entity schema.
type Project struct {
	ID          string     `json:"id"`
	ClientID    string     `json:"client_id"`
	Name        string     `json:"name"`
	Status      string     `json:"status"` // "Backlog", "Em andamento", "Go-Live"
	IsActive    bool       `json:"is_active"`
	ActivatedAt *time.Time `json:"activated_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
