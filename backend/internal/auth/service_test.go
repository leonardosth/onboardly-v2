package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashAndCheckPassword(t *testing.T) {
	password := "my_secure_password"

	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	// Check correct password
	assert.True(t, CheckPasswordHash(password, hash))

	// Check wrong password
	assert.False(t, CheckPasswordHash("wrong_password", hash))
}

func TestHashPassword_DifferentHashes(t *testing.T) {
	password := "same_password"

	hash1, err := HashPassword(password)
	require.NoError(t, err)

	hash2, err := HashPassword(password)
	require.NoError(t, err)

	// bcrypt should produce different hashes for the same password (different salt)
	assert.NotEqual(t, hash1, hash2)

	// But both should verify correctly
	assert.True(t, CheckPasswordHash(password, hash1))
	assert.True(t, CheckPasswordHash(password, hash2))
}

func TestCheckPasswordHash_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("valid_password")
	require.NoError(t, err)

	assert.False(t, CheckPasswordHash("", hash))
}

func TestGenerateJWT(t *testing.T) {
	email := "test@test.com"
	role := "Admin"
	secret := "test_secret"

	token, err := GenerateJWT(email, role, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestGenerateJWT_ValidClaims(t *testing.T) {
	email := "analyst@example.com"
	role := "Analista"
	secret := "jwt-test-secret"

	tokenStr, err := GenerateJWT(email, role, secret)
	require.NoError(t, err)

	// Parse and validate the token
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	require.NoError(t, err)
	assert.True(t, token.Valid)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)

	// Verify expiration is ~24h in the future
	assert.NotNil(t, claims.ExpiresAt)
	expiresIn := time.Until(claims.ExpiresAt.Time)
	assert.InDelta(t, 24*time.Hour, expiresIn, float64(5*time.Minute))

	// Verify issued at is recent
	assert.NotNil(t, claims.IssuedAt)
	issuedAgo := time.Since(claims.IssuedAt.Time)
	assert.Less(t, issuedAgo, 5*time.Second)
}

func TestGenerateJWT_InvalidSecret(t *testing.T) {
	tokenStr, err := GenerateJWT("user@test.com", "Admin", "correct-secret")
	require.NoError(t, err)

	// Attempt to parse with wrong secret should fail
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte("wrong-secret"), nil
	})

	assert.Error(t, err)
}

func TestCreateUser_InvalidRole(t *testing.T) {
	user, err := CreateUser("test@test.com", "password123", "InvalidRole")

	assert.Error(t, err)
	assert.Equal(t, "invalid role, must be Admin or Analista", err.Error())
	assert.Nil(t, user)
}
