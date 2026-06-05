package project

import (
	"database/sql"
	"errors"

	"onboardly-backend/internal/db"
)

// GetProjects retrieves projects filtered optionally by client_id.
func GetProjects(clientID string) ([]Project, error) {
	var rows *sql.Rows
	var err error

	if clientID != "" {
		query := `SELECT id, client_id, name, status, is_active, created_at, updated_at FROM projects WHERE client_id = $1 ORDER BY created_at DESC`
		rows, err = db.DB.Query(query, clientID)
	} else {
		query := `SELECT id, client_id, name, status, is_active, created_at, updated_at FROM projects ORDER BY created_at DESC`
		rows, err = db.DB.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.ClientID, &p.Name, &p.Status, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

// GetProjectByID retrieves a single project by UUID.
func GetProjectByID(id string) (*Project, error) {
	query := `SELECT id, client_id, name, status, is_active, created_at, updated_at FROM projects WHERE id = $1`
	var p Project
	err := db.DB.QueryRow(query, id).Scan(&p.ID, &p.ClientID, &p.Name, &p.Status, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return &p, nil
}

// CreateProject inserts a new project record.
func CreateProject(clientID, name string) (*Project, error) {
	if name == "" {
		return nil, errors.New("project name cannot be empty")
	}

	// Verify client exists
	var dummy string
	clientCheckQuery := `SELECT id FROM clients WHERE id = $1`
	err := db.DB.QueryRow(clientCheckQuery, clientID).Scan(&dummy)
	if err != nil {
		return nil, errors.New("associated client not found")
	}

	query := `INSERT INTO projects (client_id, name, status, is_active) VALUES ($1, $2, 'Backlog', TRUE) RETURNING id, status, is_active, created_at, updated_at`
	p := &Project{ClientID: clientID, Name: name}
	err = db.DB.QueryRow(query, p.ClientID, p.Name).Scan(&p.ID, &p.Status, &p.IsActive, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// UpdateProjectStatus handles transitioning status stages and active state logic.
func UpdateProjectStatus(id, status string) (*Project, error) {
	if status != "Backlog" && status != "Em andamento" && status != "Go-Live" {
		return nil, errors.New("invalid status: must be Backlog, Em andamento, or Go-Live")
	}

	// is_active is false only when status is Go-Live
	isActive := status != "Go-Live"

	query := `UPDATE projects SET status = $1, is_active = $2, updated_at = NOW() WHERE id = $3 RETURNING client_id, name, created_at, updated_at`
	p := &Project{ID: id, Status: status, IsActive: isActive}
	err := db.DB.QueryRow(query, p.Status, p.IsActive, p.ID).Scan(&p.ClientID, &p.Name, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return p, nil
}
