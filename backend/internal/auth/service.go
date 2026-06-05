package auth

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"onboardly-backend/internal/db"
)

// Claims represents the JWT claims schema.
type Claims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword generates a bcrypt hash of a plain text password.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with a hashed password.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a signed JWT token valid for 24 hours.
func GenerateJWT(email, role, secret string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// AuthenticateUser verifies the login credentials and returns a JWT if valid.
func AuthenticateUser(email, password, secret string) (string, string, error) {
	query := `SELECT password_hash, role FROM users WHERE email = $1`
	var hash, role string
	err := db.DB.QueryRow(query, email).Scan(&hash, &role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", "", errors.New("invalid email or password")
		}
		return "", "", err
	}

	if !CheckPasswordHash(password, hash) {
		return "", "", errors.New("invalid email or password")
	}

	token, err := GenerateJWT(email, role, secret)
	if err != nil {
		return "", "", err
	}

	return token, role, nil
}

// CreateUser inserts a new user record into the database.
func CreateUser(email, password, role string) (*User, error) {
	if role != "Admin" && role != "Analista" {
		return nil, errors.New("invalid role, must be Admin or Analista")
	}

	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	query := `INSERT INTO users (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id, created_at`
	var id string
	var createdAt time.Time
	err = db.DB.QueryRow(query, email, hash, role).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        id,
		Email:     email,
		Role:      role,
		CreatedAt: createdAt,
	}, nil
}
