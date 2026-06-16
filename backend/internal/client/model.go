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

// ClientWithDetails includes aggregated data about the client's latest project and meetings.
type ClientWithDetails struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	CNPJ              string     `json:"cnpj"`
	ProjectName       *string    `json:"project_name"`
	ProjectStatus     *string    `json:"project_status"`
	ProjectIsActive   *bool      `json:"project_is_active"`
	Responsible       *string    `json:"responsible"`
	CompletedAgendas  int        `json:"completed_agendas"`
	TotalAgendas      int        `json:"total_agendas"`
	CreatedAt         time.Time  `json:"created_at"`
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
