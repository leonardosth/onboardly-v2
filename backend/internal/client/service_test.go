package client

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"onboardly-backend/internal/db"
)

func TestGetClientsWithDetails(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer mockDB.Close()

	// Replace the global DB with the mock
	originalDB := db.DB
	db.DB = mockDB
	defer func() { db.DB = originalDB }()

	now := time.Now()

	t.Run("returns clients with details", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "name", "cnpj", "created_at",
			"project_name", "project_status", "is_active",
			"responsible", "completed_agendas", "total_agendas",
		}).
			AddRow("client-1", "Empresa A", "11.111.111/0001-11", now,
				"Projeto Alpha", "Em andamento", true,
				"analista@test.com", 2, 5).
			AddRow("client-2", "Empresa B", "22.222.222/0001-22", now,
				nil, nil, nil,
				nil, 0, 0)

		mock.ExpectQuery("SELECT c.id, c.name, c.cnpj, c.created_at").WillReturnRows(rows)

		clients, err := GetClientsWithDetails()

		assert.NoError(t, err)
		assert.Len(t, clients, 2)

		assert.Equal(t, "Empresa A", clients[0].Name)
		assert.Equal(t, "Projeto Alpha", *clients[0].ProjectName)
		assert.Equal(t, "Em andamento", *clients[0].ProjectStatus)
		assert.True(t, *clients[0].ProjectIsActive)
		assert.Equal(t, "analista@test.com", *clients[0].Responsible)
		assert.Equal(t, 2, clients[0].CompletedAgendas)

		assert.Equal(t, "Empresa B", clients[1].Name)
		assert.Nil(t, clients[1].ProjectName)
		assert.Equal(t, 0, clients[1].TotalAgendas)
	})

	t.Run("handles db error", func(t *testing.T) {
		mock.ExpectQuery("SELECT c.id, c.name, c.cnpj, c.created_at").WillReturnError(sql.ErrConnDone)

		clients, err := GetClientsWithDetails()

		assert.Error(t, err)
		assert.Nil(t, clients)
	})
}
