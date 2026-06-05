package meeting

import (
	"errors"
	"time"

	"onboardly-backend/internal/db"
)

// GetMeetingsByProject retrieves all scheduled meetings for a project.
func GetMeetingsByProject(projectID string) ([]Meeting, error) {
	query := `SELECT id, project_id, analyst_id, title, scheduled_at, no_show, created_at FROM meetings WHERE project_id = $1 ORDER BY scheduled_at ASC`
	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []Meeting
	for rows.Next() {
		var m Meeting
		var analystID sqlNullString
		err := rows.Scan(&m.ID, &m.ProjectID, &analystID, &m.Title, &m.ScheduledAt, &m.NoShow, &m.CreatedAt)
		if err != nil {
			return nil, err
		}
		if analystID.Valid {
			m.AnalystID = analystID.String
		}
		meetings = append(meetings, m)
	}

	return meetings, nil
}

// CreateMeeting inserts a meeting and enforces project existence constraints.
func CreateMeeting(projectID, analystID, title string, scheduledAt time.Time) (*Meeting, error) {
	if title == "" {
		return nil, errors.New("meeting title cannot be empty")
	}

	// Verify project exists
	var dummy string
	projectCheckQuery := `SELECT id FROM projects WHERE id = $1`
	err := db.DB.QueryRow(projectCheckQuery, projectID).Scan(&dummy)
	if err != nil {
		return nil, errors.New("cannot schedule meeting: associated project not found")
	}

	var analystVal interface{}
	if analystID != "" {
		analystVal = analystID
	} else {
		analystVal = nil
	}

	query := `INSERT INTO meetings (project_id, analyst_id, title, scheduled_at) VALUES ($1, $2, $3, $4) RETURNING id, no_show, created_at`
	m := &Meeting{
		ProjectID:   projectID,
		AnalystID:   analystID,
		Title:       title,
		ScheduledAt: scheduledAt,
	}

	err = db.DB.QueryRow(query, m.ProjectID, analystVal, m.Title, m.ScheduledAt).Scan(&m.ID, &m.NoShow, &m.CreatedAt)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// sqlNullString is a helper for scanning nullable strings.
type sqlNullString struct {
	String string
	Valid  bool
}

func (s *sqlNullString) Scan(value interface{}) error {
	if value == nil {
		s.String, s.Valid = "", false
		return nil
	}
	s.Valid = true
	switch v := value.(type) {
	case string:
		s.String = v
	case []byte:
		s.String = string(v)
	}
	return nil
}
