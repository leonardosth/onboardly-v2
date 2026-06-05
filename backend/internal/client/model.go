package client

import (
	"errors"
	"regexp"
	"time"
)

// Client represents the client entity schema.
type Client struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CNPJ      string    `json:"cnpj"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var cnpjRegex = regexp.MustCompile(`^\d{2}\.\d{3}\.\d{3}/\d{4}-\d{2}$`)

// Validate checks client attribute rules.
func (c *Client) Validate() error {
	if c.Name == "" {
		return errors.New("client name cannot be empty")
	}
	if !cnpjRegex.MatchString(c.CNPJ) {
		return errors.New("invalid CNPJ format, must be XX.XXX.XXX/XXXX-XX")
	}
	return nil
}
