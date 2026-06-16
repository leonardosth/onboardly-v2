package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Validate(t *testing.T) {
	tests := []struct {
		name    string
		client  Client
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid client",
			client: Client{
				Name: "Empresa XYZ",
				CNPJ: "12.345.678/0001-90",
			},
			wantErr: false,
		},
		{
			name: "empty name",
			client: Client{
				Name: "",
				CNPJ: "12.345.678/0001-90",
			},
			wantErr: true,
			errMsg:  "client name cannot be empty",
		},
		{
			name: "invalid CNPJ format",
			client: Client{
				Name: "Empresa XYZ",
				CNPJ: "12345678000190",
			},
			wantErr: true,
			errMsg:  "invalid CNPJ format, must be XX.XXX.XXX/XXXX-XX",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.client.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
