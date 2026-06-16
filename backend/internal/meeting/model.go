package meeting

import "time"

// Meeting represents the meeting entity schema.
type Meeting struct {
	ID          string     `json:"id"`
	ProjectID   string     `json:"project_id"`
	AnalystID   string     `json:"analyst_id"`
	Title       string     `json:"title"`
	ScheduledAt time.Time  `json:"scheduled_at"`
	Status      string     `json:"status"`
	CompletedAt *time.Time `json:"completed_at"`
	NoShow      bool       `json:"no_show"`
	CreatedAt   time.Time  `json:"created_at"`
}

// MeetingWithDetails includes joined project and client names for display.
type MeetingWithDetails struct {
	Meeting
	ProjectName string `json:"project_name"`
	ClientName  string `json:"client_name"`
}
