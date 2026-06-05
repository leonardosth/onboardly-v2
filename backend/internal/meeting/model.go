package meeting

import "time"

// Meeting represents the meeting entity schema.
type Meeting struct {
	ID          string    `json:"id"`
	ProjectID   string    `json:"project_id"`
	AnalystID   string    `json:"analyst_id"`
	Title       string    `json:"title"`
	ScheduledAt time.Time `json:"scheduled_at"`
	NoShow      bool      `json:"no_show"`
	CreatedAt   time.Time `json:"created_at"`
}
