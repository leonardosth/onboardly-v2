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
