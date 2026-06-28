package dashboard

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

func TestGetDashboardData(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	t.Run("returns dashboard data with metrics", func(t *testing.T) {
		// 1. Total projects
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
		// 2. Go-Live projects
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
		// 3. Total meetings
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(20))
		// 4. No-show meetings
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		// 5. Funnel: Registered (total projects)
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
		// 6. Funnel: Participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))
		// 7. Funnel: Active
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))
		// 8. Abandonment: total projects
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
		// 9. Abandonment: abandoned count
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4))
		// 10. 30d activation: participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))
		// 11. 30d activation: activated in 30d
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(5))
		// 12. First meeting activation: total participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(7))
		// 13. First meeting activation: first meeting activated
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		// 14. Cohorts
		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"cohort", "total", "activated"}).
				AddRow("2026-06", 5, 2).
				AddRow("2026-05", 3, 1))
		// 15. History
		mock.ExpectQuery("SELECT TO_CHAR").
			WillReturnRows(sqlmock.NewRows([]string{"month", "count"}).
				AddRow("2026-06", 2).
				AddRow("2026-05", 1))
		// 16. Activities
		mock.ExpectQuery("SELECT id, entity_type, entity_id, description, created_at FROM activities").
			WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "description", "created_at"}).
				AddRow("a-1", "client", "c-1", "Client created", now))

		data, err := GetDashboardData()

		assert.NoError(t, err)
		require.NotNil(t, data)

		// Activation rate: 3/10 * 100 = 30.0
		assert.Equal(t, 30.0, data.Metrics.ActivationRate)
		// No-show rate: 2/20 * 100 = 10.0
		assert.Equal(t, 10.0, data.Metrics.NoShowRate)
		// Funnel
		assert.Equal(t, 10, data.Funnel.Registered)
		assert.Equal(t, 7, data.Funnel.Participants)
		assert.Equal(t, 3, data.Funnel.Active)
		// History
		assert.Len(t, data.History, 2)
		// Cohorts
		assert.Len(t, data.Cohorts, 2)
		// Activities
		assert.Len(t, data.RecentActivities, 1)
	})
}

func TestGetDashboardData_ZeroDivision(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	t.Run("handles zero projects gracefully", func(t *testing.T) {
		// Total projects = 0
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Go-Live = 0
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Total meetings = 0
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// No-show = 0
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Funnel: registered
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Funnel: participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Funnel: active
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Abandonment: total
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Abandonment: abandoned
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// 30d: participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// 30d: activated
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// First meeting: participants
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// First meeting: activated
		mock.ExpectQuery("SELECT COUNT").
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		// Cohorts
		mock.ExpectQuery("SELECT").
			WillReturnRows(sqlmock.NewRows([]string{"cohort", "total", "activated"}))
		// History
		mock.ExpectQuery("SELECT TO_CHAR").
			WillReturnRows(sqlmock.NewRows([]string{"month", "count"}))
		// Activities
		mock.ExpectQuery("SELECT id, entity_type, entity_id, description, created_at FROM activities").
			WillReturnRows(sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "description", "created_at"}))

		data, err := GetDashboardData()

		assert.NoError(t, err)
		require.NotNil(t, data)

		// All rates should be 0 (no division by zero panic)
		assert.Equal(t, 0.0, data.Metrics.ActivationRate)
		assert.Equal(t, 0.0, data.Metrics.NoShowRate)
		assert.Equal(t, 0.0, data.Metrics.AbandonmentRate)
		assert.Equal(t, 0.0, data.Metrics.FirstMeetingActivationRate)
		assert.Equal(t, 0.0, data.Metrics.Activation30dRate)

		// History should have fallback entries when empty
		assert.Len(t, data.History, 6)
	})
}
