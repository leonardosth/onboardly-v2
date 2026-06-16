package meeting

import (
	"errors"
	"time"

	"onboardly-backend/internal/db"
)

// GetMeetingsByProject retrieves all meetings for a project.
func GetMeetingsByProject(projectID string) ([]Meeting, error) {
	query := `SELECT id, project_id, analyst_id, title, scheduled_at, status, completed_at, no_show, created_at FROM meetings WHERE project_id = $1 ORDER BY scheduled_at ASC`
	rows, err := db.DB.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []Meeting
	for rows.Next() {
		var m Meeting
		var analystID sqlNullString
		err := rows.Scan(&m.ID, &m.ProjectID, &analystID, &m.Title, &m.ScheduledAt, &m.Status, &m.CompletedAt, &m.NoShow, &m.CreatedAt)
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

// GetMeetingsByAnalyst retrieves meetings for a specific analyst, optionally filtered by status.
func GetMeetingsByAnalyst(analystID, status string) ([]MeetingWithDetails, error) {
	query := `
		SELECT m.id, m.project_id, m.analyst_id, m.title, m.scheduled_at, m.status, m.completed_at, m.no_show, m.created_at,
		       p.name AS project_name, c.name AS client_name
		FROM meetings m
		JOIN projects p ON p.id = m.project_id
		JOIN clients c ON c.id = p.client_id
		WHERE m.analyst_id = $1
	`
	args := []interface{}{analystID}

	if status != "" {
		query += ` AND m.status = $2`
		args = append(args, status)
	}

	query += ` ORDER BY m.scheduled_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meetings []MeetingWithDetails
	for rows.Next() {
		var md MeetingWithDetails
		var aID sqlNullString
		err := rows.Scan(
			&md.ID, &md.ProjectID, &aID, &md.Title, &md.ScheduledAt, &md.Status, &md.CompletedAt, &md.NoShow, &md.CreatedAt,
			&md.ProjectName, &md.ClientName,
		)
		if err != nil {
			return nil, err
		}
		if aID.Valid {
			md.AnalystID = aID.String
		}
		meetings = append(meetings, md)
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

	query := `INSERT INTO meetings (project_id, analyst_id, title, scheduled_at, status) VALUES ($1, $2, $3, $4, 'scheduled') RETURNING id, status, completed_at, no_show, created_at`
	m := &Meeting{
		ProjectID:   projectID,
		AnalystID:   analystID,
		Title:       title,
		ScheduledAt: scheduledAt,
	}

	err = db.DB.QueryRow(query, m.ProjectID, analystVal, m.Title, m.ScheduledAt).Scan(&m.ID, &m.Status, &m.CompletedAt, &m.NoShow, &m.CreatedAt)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// CompleteMeeting marks a meeting as completed and optionally activates the client.
func CompleteMeeting(meetingID string, activateClient bool) (*Meeting, bool, error) {
	// Fetch meeting
	var m Meeting
	var analystID sqlNullString
	fetchQuery := `SELECT id, project_id, analyst_id, title, scheduled_at, status, completed_at, no_show, created_at FROM meetings WHERE id = $1`
	err := db.DB.QueryRow(fetchQuery, meetingID).Scan(&m.ID, &m.ProjectID, &analystID, &m.Title, &m.ScheduledAt, &m.Status, &m.CompletedAt, &m.NoShow, &m.CreatedAt)
	if err != nil {
		return nil, false, errors.New("meeting not found")
	}
	if analystID.Valid {
		m.AnalystID = analystID.String
	}

	if m.Status != "scheduled" {
		return nil, false, errors.New("meeting is already " + m.Status)
	}

	// Update meeting to completed
	now := time.Now()
	updateQuery := `UPDATE meetings SET status = 'completed', completed_at = $1 WHERE id = $2`
	_, err = db.DB.Exec(updateQuery, now, meetingID)
	if err != nil {
		return nil, false, err
	}
	m.Status = "completed"
	m.CompletedAt = &now

	// Optionally activate client
	projectActivated := false
	if activateClient {
		activateQuery := `UPDATE projects SET activated_at = $1, updated_at = $1 WHERE id = $2 AND activated_at IS NULL`
		res, err := db.DB.Exec(activateQuery, now, m.ProjectID)
		if err != nil {
			return &m, false, nil // Meeting completed but activation failed silently
		}
		rows, _ := res.RowsAffected()
		projectActivated = rows > 0
	}

	return &m, projectActivated, nil
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
