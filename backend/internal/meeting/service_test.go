package meeting

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

func TestCreateMeeting_EmptyTitle(t *testing.T) {
	_, _, cleanup := setupMockDB(t)
	defer cleanup()

	m, err := CreateMeeting("project-1", "analyst-1", "", time.Now())

	assert.Error(t, err)
	assert.Equal(t, "meeting title cannot be empty", err.Error())
	assert.Nil(t, m)
}

func TestCreateMeeting_ProjectNotFound(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id FROM projects WHERE id").
		WithArgs("nonexistent-project").
		WillReturnError(sql.ErrNoRows)

	m, err := CreateMeeting("nonexistent-project", "analyst-1", "Meeting Title", time.Now())

	assert.Error(t, err)
	assert.Equal(t, "cannot schedule meeting: associated project not found", err.Error())
	assert.Nil(t, m)
}

func TestGetMeetingsByProject(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	t.Run("returns meetings for project", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "project_id", "analyst_id", "title", "scheduled_at", "status", "completed_at", "no_show", "created_at",
		}).
			AddRow("m-1", "p-1", "a-1", "Kickoff", now, "scheduled", nil, false, now).
			AddRow("m-2", "p-1", nil, "Follow-up", now, "completed", now, false, now)

		mock.ExpectQuery("SELECT id, project_id, analyst_id, title, scheduled_at, status, completed_at, no_show, created_at FROM meetings WHERE project_id").
			WithArgs("p-1").
			WillReturnRows(rows)

		meetings, err := GetMeetingsByProject("p-1")

		assert.NoError(t, err)
		assert.Len(t, meetings, 2)
		assert.Equal(t, "Kickoff", meetings[0].Title)
		assert.Equal(t, "a-1", meetings[0].AnalystID)
		assert.Equal(t, "Follow-up", meetings[1].Title)
		assert.Empty(t, meetings[1].AnalystID)
	})

	t.Run("returns error on db failure", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, project_id, analyst_id, title").
			WithArgs("p-1").
			WillReturnError(sql.ErrConnDone)

		meetings, err := GetMeetingsByProject("p-1")

		assert.Error(t, err)
		assert.Nil(t, meetings)
	})
}

func TestGetMeetingsByAnalyst(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()

	t.Run("returns meetings for analyst", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "project_id", "analyst_id", "title", "scheduled_at", "status", "completed_at", "no_show", "created_at",
			"project_name", "client_name",
		}).
			AddRow("m-1", "p-1", "a-1", "Meeting 1", now, "scheduled", nil, false, now, "Projeto Alpha", "Empresa A")

		mock.ExpectQuery("SELECT m.id, m.project_id, m.analyst_id").
			WithArgs("a-1").
			WillReturnRows(rows)

		meetings, err := GetMeetingsByAnalyst("a-1", "")

		assert.NoError(t, err)
		assert.Len(t, meetings, 1)
		assert.Equal(t, "Meeting 1", meetings[0].Title)
		assert.Equal(t, "Projeto Alpha", meetings[0].ProjectName)
		assert.Equal(t, "Empresa A", meetings[0].ClientName)
	})

	t.Run("filters by status", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{
			"id", "project_id", "analyst_id", "title", "scheduled_at", "status", "completed_at", "no_show", "created_at",
			"project_name", "client_name",
		}).
			AddRow("m-2", "p-1", "a-1", "Completed Meeting", now, "completed", now, false, now, "Projeto Beta", "Empresa B")

		mock.ExpectQuery("SELECT m.id, m.project_id, m.analyst_id").
			WithArgs("a-1", "completed").
			WillReturnRows(rows)

		meetings, err := GetMeetingsByAnalyst("a-1", "completed")

		assert.NoError(t, err)
		assert.Len(t, meetings, 1)
		assert.Equal(t, "completed", meetings[0].Status)
	})
}

func TestCompleteMeeting_AlreadyCompleted(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	now := time.Now()
	completedAt := now

	rows := sqlmock.NewRows([]string{
		"id", "project_id", "analyst_id", "title", "scheduled_at", "status", "completed_at", "no_show", "created_at",
	}).AddRow("m-1", "p-1", "a-1", "Old Meeting", now, "completed", completedAt, false, now)

	mock.ExpectQuery("SELECT id, project_id, analyst_id, title, scheduled_at, status, completed_at, no_show, created_at FROM meetings WHERE id").
		WithArgs("m-1").
		WillReturnRows(rows)

	m, activated, err := CompleteMeeting("m-1", false)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "meeting is already completed")
	assert.Nil(t, m)
	assert.False(t, activated)
}

func TestCompleteMeeting_NotFound(t *testing.T) {
	_, mock, cleanup := setupMockDB(t)
	defer cleanup()

	mock.ExpectQuery("SELECT id, project_id, analyst_id, title, scheduled_at, status, completed_at, no_show, created_at FROM meetings WHERE id").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	m, activated, err := CompleteMeeting("nonexistent", false)

	assert.Error(t, err)
	assert.Equal(t, "meeting not found", err.Error())
	assert.Nil(t, m)
	assert.False(t, activated)
}
