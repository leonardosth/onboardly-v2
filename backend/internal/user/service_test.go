package user

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

func TestDeleteUser_SelfDeletion(t *testing.T) {
	err := DeleteUser("user-123", "user-123")

	assert.Error(t, err)
	assert.Equal(t, "cannot delete your own account", err.Error())
}

func TestDeleteUser_NotFound(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT role FROM users WHERE id").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	err := DeleteUser("nonexistent", "other-user")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestDeleteUser_LastAdmin(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// Return "Admin" role for the user being deleted
	mock.ExpectQuery("SELECT role FROM users WHERE id").
		WithArgs("admin-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("Admin"))

	// Return admin count = 1 (last admin)
	mock.ExpectQuery("SELECT count").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	err := DeleteUser("admin-1", "other-user")

	assert.Error(t, err)
	assert.Equal(t, "cannot delete the last admin in the system", err.Error())
}

func TestDeleteUser_Success(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	// Return "Analista" role
	mock.ExpectQuery("SELECT role FROM users WHERE id").
		WithArgs("analyst-1").
		WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("Analista"))

	// Expect delete
	mock.ExpectExec("DELETE FROM users WHERE id").
		WithArgs("analyst-1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := DeleteUser("analyst-1", "admin-1")

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListUsers(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	t.Run("returns users list", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"}).
			AddRow("u-1", "admin@test.com", "Admin", now).
			AddRow("u-2", "analyst@test.com", "Analista", now)

		mock.ExpectQuery("SELECT id, email, role, created_at FROM users").
			WillReturnRows(rows)

		users, err := ListUsers()

		assert.NoError(t, err)
		assert.Len(t, users, 2)
		assert.Equal(t, "admin@test.com", users[0].Email)
		assert.Equal(t, "Admin", users[0].Role)
		assert.Equal(t, "analyst@test.com", users[1].Email)
	})

	t.Run("returns nil when no users", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "email", "role", "created_at"})
		mock.ExpectQuery("SELECT id, email, role, created_at FROM users").
			WillReturnRows(rows)

		users, err := ListUsers()

		assert.NoError(t, err)
		assert.Nil(t, users)
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, role, created_at FROM users").
			WillReturnError(sql.ErrConnDone)

		users, err := ListUsers()

		assert.Error(t, err)
		assert.Nil(t, users)
	})
}
