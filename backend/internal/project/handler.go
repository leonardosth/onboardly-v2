package project

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"onboardly-backend/internal/apierr"
)

type projectRequest struct {
	ClientID string `json:"client_id"`
	Name     string `json:"name"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

// ListProjectsHandler returns a list of projects.
func ListProjectsHandler(w http.ResponseWriter, r *http.Request) {
	clientID := r.URL.Query().Get("client_id")
	projects, err := GetProjects(clientID)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(projects)
}

// GetProjectHandler returns project details.
func GetProjectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	p, err := GetProjectByID(id)
	if err != nil {
		apierr.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(p)
}

// CreateProjectHandler creates a new project.
func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var req projectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := CreateProject(req.ClientID, req.Name)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(p)
}

// UpdateProjectStatusHandler handles status transition request.
func UpdateProjectStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := UpdateProjectStatus(id, req.Status)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(p)
}

// FinalizeProjectHandler sets a project status to Go-Live (finalized).
func FinalizeProjectHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	p, err := FinalizeProject(id)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(p)
}
