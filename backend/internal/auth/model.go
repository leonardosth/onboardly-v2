package auth

import "time"

// User represents the database entity for a platform operator.
type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"` // "Admin" or "Analista"
	CreatedAt    time.Time `json:"created_at"`
}
