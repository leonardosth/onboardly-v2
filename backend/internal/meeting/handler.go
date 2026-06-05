package meeting

import (
	"encoding/json"
	"net/http"
	"time"

	"onboardly-backend/internal/apierr"
	"onboardly-backend/internal/auth"
	"onboardly-backend/internal/db"
)

type meetingRequest struct {
	ProjectID   string    `json:"project_id"`
	AnalystID   string    `json:"analyst_id"`
	Title       string    `json:"title"`
	ScheduledAt time.Time `json:"scheduled_at"` // Expects ISO 8601
}

// ListMeetingsHandler lists scheduled meetings for a project.
func ListMeetingsHandler(w http.ResponseWriter, r *http.Request) {
	projectID := r.URL.Query().Get("project_id")
	if projectID == "" {
		apierr.WriteError(w, http.StatusBadRequest, "query parameter project_id is required")
		return
	}

	meetings, err := GetMeetingsByProject(projectID)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(meetings)
}

// CreateMeetingHandler schedules a new meeting.
func CreateMeetingHandler(w http.ResponseWriter, r *http.Request) {
	var req meetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ProjectID == "" || req.Title == "" || req.ScheduledAt.IsZero() {
		apierr.WriteError(w, http.StatusBadRequest, "project_id, title, and scheduled_at are required")
		return
	}

	emailVal := r.Context().Value(auth.UserEmailKey)
	if emailVal == nil {
		apierr.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var actualAnalystID string
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", emailVal.(string)).Scan(&actualAnalystID)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, "failed to get user identity")
		return
	}

	m, err := CreateMeeting(req.ProjectID, actualAnalystID, req.Title, req.ScheduledAt)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(m)
}
