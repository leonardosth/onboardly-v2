package activity

import (
	"time"

	"onboardly-backend/internal/db"
)

// Activity represents the structure of an activity log.
type Activity struct {
	ID          string    `json:"id"`
	EntityType  string    `json:"entity_type"`
	EntityID    string    `json:"entity_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// LogActivity inserts a new audit log record.
func LogActivity(entityType, entityID, description string) error {
	query := `INSERT INTO activities (entity_type, entity_id, description) VALUES ($1, $2, $3)`
	_, err := db.DB.Exec(query, entityType, entityID, description)
	return err
}

// GetRecentActivities retrieves the 15 most recent activities.
func GetRecentActivities() ([]Activity, error) {
	query := `SELECT id, entity_type, entity_id, description, created_at FROM activities ORDER BY created_at DESC LIMIT 15`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []Activity
	for rows.Next() {
		var a Activity
		err := rows.Scan(&a.ID, &a.EntityType, &a.EntityID, &a.Description, &a.CreatedAt)
		if err != nil {
			return nil, err
		}
		logs = append(logs, a)
	}

	return logs, nil
}
