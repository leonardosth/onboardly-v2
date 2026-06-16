package client

import (
	"database/sql"
	"errors"

	"onboardly-backend/internal/db"
)

// GetAllClients returns a list of all client records.
func GetAllClients() ([]Client, error) {
	query := `SELECT id, name, cnpj, created_at, updated_at FROM clients ORDER BY name ASC`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var c Client
		err := rows.Scan(&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return clients, nil
}

// GetClientsWithDetails returns a list of clients including aggregated project and meeting details.
func GetClientsWithDetails() ([]ClientWithDetails, error) {
	query := `
		SELECT c.id, c.name, c.cnpj, c.created_at,
		       p.name AS project_name, p.status AS project_status, p.is_active,
		       u.email AS responsible,
		       COUNT(m.id) FILTER (WHERE m.status = 'completed') AS completed_agendas,
		       COUNT(m.id) AS total_agendas
		FROM clients c
		LEFT JOIN LATERAL (
		    SELECT * FROM projects
		    WHERE client_id = c.id
		    ORDER BY created_at DESC LIMIT 1
		) p ON TRUE
		LEFT JOIN meetings m ON m.project_id = p.id
		LEFT JOIN LATERAL (
		    SELECT analyst_id FROM meetings
		    WHERE project_id = p.id AND analyst_id IS NOT NULL
		    ORDER BY scheduled_at DESC LIMIT 1
		) latest_m ON TRUE
		LEFT JOIN users u ON u.id = latest_m.analyst_id
		GROUP BY c.id, c.name, c.cnpj, c.created_at, p.name, p.status, p.is_active, u.email
		ORDER BY c.name ASC
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []ClientWithDetails
	for rows.Next() {
		var c ClientWithDetails
		err := rows.Scan(
			&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt,
			&c.ProjectName, &c.ProjectStatus, &c.ProjectIsActive,
			&c.Responsible, &c.CompletedAgendas, &c.TotalAgendas,
		)
		if err != nil {
			return nil, err
		}
		clients = append(clients, c)
	}

	return clients, nil
}

// GetClientByID retrieves a single client by UUID.
func GetClientByID(id string) (*Client, error) {
	query := `SELECT id, name, cnpj, created_at, updated_at FROM clients WHERE id = $1`
	var c Client
	err := db.DB.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &c, nil
}

// GetClientByCNPJ retrieves a client by CNPJ.
func GetClientByCNPJ(cnpj string) (*Client, error) {
	query := `SELECT id, name, cnpj, created_at, updated_at FROM clients WHERE cnpj = $1`
	var c Client
	err := db.DB.QueryRow(query, cnpj).Scan(&c.ID, &c.Name, &c.CNPJ, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("client not found")
		}
		return nil, err
	}
	return &c, nil
}

// InsertClient adds a new client record into the database.
func InsertClient(name, cnpj string) (*Client, error) {
	c := &Client{Name: name, CNPJ: cnpj}
	if err := c.Validate(); err != nil {
		return nil, err
	}

	query := `INSERT INTO clients (name, cnpj) VALUES ($1, $2) RETURNING id, created_at, updated_at`
	err := db.DB.QueryRow(query, c.Name, c.CNPJ).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// UpdateClientDetails modifies client fields.
func UpdateClientDetails(id, name, cnpj string) (*Client, error) {
	c := &Client{ID: id, Name: name, CNPJ: cnpj}
	if err := c.Validate(); err != nil {
		return nil, err
	}

	query := `UPDATE clients SET name = $1, cnpj = $2, updated_at = NOW() WHERE id = $3 RETURNING created_at, updated_at`
	err := db.DB.QueryRow(query, c.Name, c.CNPJ, c.ID).Scan(&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// RemoveClient deletes a client record.
func RemoveClient(id string) error {
	query := `DELETE FROM clients WHERE id = $1`
	res, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("client not found")
	}

	return nil
}
