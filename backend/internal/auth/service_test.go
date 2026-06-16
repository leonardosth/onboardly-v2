package auth

import (
	"testing"

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

func TestGenerateJWT(t *testing.T) {
	email := "test@test.com"
	role := "Admin"
	secret := "test_secret"

	token, err := GenerateJWT(email, role, secret)
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}
