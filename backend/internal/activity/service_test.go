package activity

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"onboardly-backend/internal/db"
)

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	require.NoError(t, err)

	originalDB := db.DB
	db.DB = mockDB

	cleanup := func() {
		db.DB = originalDB
		mockDB.Close()
	}

	return mockDB, mock, cleanup
}

func TestLogActivity(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("inserts activity successfully", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO activities").
			WithArgs("client", "client-123", "Client created").
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := LogActivity("client", "client-123", "Client created")

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		mock.ExpectExec("INSERT INTO activities").
			WithArgs("project", "proj-1", "Status changed").
			WillReturnError(sql.ErrConnDone)

		err := LogActivity("project", "proj-1", "Status changed")

		assert.Error(t, err)
	})
}

func TestGetRecentActivities(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	t.Run("returns recent activities", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "description", "created_at"}).
			AddRow("act-1", "client", "c-1", "Client created", now).
			AddRow("act-2", "project", "p-1", "Project started", now)

		mock.ExpectQuery("SELECT id, entity_type, entity_id, description, created_at FROM activities").
			WillReturnRows(rows)

		activities, err := GetRecentActivities()

		assert.NoError(t, err)
		assert.Len(t, activities, 2)
		assert.Equal(t, "Client created", activities[0].Description)
		assert.Equal(t, "client", activities[0].EntityType)
		assert.Equal(t, "Project started", activities[1].Description)
	})

	t.Run("returns empty list when no activities", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "description", "created_at"})
		mock.ExpectQuery("SELECT id, entity_type, entity_id, description, created_at FROM activities").
			WillReturnRows(rows)

		activities, err := GetRecentActivities()

		assert.NoError(t, err)
		assert.Nil(t, activities)
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, entity_type, entity_id, description, created_at FROM activities").
			WillReturnError(sql.ErrConnDone)

		activities, err := GetRecentActivities()

		assert.Error(t, err)
		assert.Nil(t, activities)
	})
}
