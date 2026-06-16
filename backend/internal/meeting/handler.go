package meeting

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
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

type completeMeetingRequest struct {
	ActivateClient bool `json:"activate_client"`
}

type completeMeetingResponse struct {
	Meeting          *Meeting `json:"meeting"`
	ProjectActivated bool     `json:"project_activated"`
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

// ListMyMeetingsHandler lists the authenticated analyst's meetings.
func ListMyMeetingsHandler(w http.ResponseWriter, r *http.Request) {
	emailVal := r.Context().Value(auth.UserEmailKey)
	if emailVal == nil {
		apierr.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Look up user ID from email
	var analystID string
	err := db.DB.QueryRow("SELECT id FROM users WHERE email = $1", emailVal.(string)).Scan(&analystID)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, "failed to get user identity")
		return
	}

	status := r.URL.Query().Get("status")
	if status == "" {
		status = "completed"
	}

	meetings, err := GetMeetingsByAnalyst(analystID, status)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if meetings == nil {
		meetings = []MeetingWithDetails{}
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

// CompleteMeetingHandler marks a meeting as completed and optionally activates the client.
func CompleteMeetingHandler(w http.ResponseWriter, r *http.Request) {
	meetingID := chi.URLParam(r, "id")
	if meetingID == "" {
		apierr.WriteError(w, http.StatusBadRequest, "meeting id is required")
		return
	}

	var req completeMeetingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// Allow empty body (activate_client defaults to false)
		req.ActivateClient = false
	}

	m, activated, err := CompleteMeeting(meetingID, req.ActivateClient)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	resp := completeMeetingResponse{
		Meeting:          m,
		ProjectActivated: activated,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
