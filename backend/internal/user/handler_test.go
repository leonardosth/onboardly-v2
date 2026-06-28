package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{"valid password", "Secure1234", true},
		{"minimum valid", "abcdefg1", true},
		{"just numbers and letters", "a1b2c3d4", true},
		{"too short", "Ab1", false},
		{"exactly 7 chars", "Abcdef1", false},
		{"no numbers", "abcdefgh", false},
		{"no letters", "12345678", false},
		{"empty string", "", false},
		{"only special chars", "!@#$%^&*", false},
		{"special chars with number no letter", "!@#$%^1*", false},
		{"special chars with letter no number", "!@#$%^a*", false},
		{"valid with special chars", "Passw0rd!", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validatePassword(tt.password)
			assert.Equal(t, tt.valid, result)
		})
	}
}
