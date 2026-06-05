package client

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"onboardly-backend/internal/apierr"
)

type clientRequest struct {
	Name string `json:"name"`
	CNPJ string `json:"cnpj"`
}

// ListClientsHandler responds with all clients list.
func ListClientsHandler(w http.ResponseWriter, r *http.Request) {
	clients, err := GetAllClients()
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(clients)
}

// GetClientHandler responds with details of a single client.
func GetClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	c, err := GetClientByID(id)
	if err != nil {
		apierr.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(c)
}

// CreateClientHandler creates a new client.
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := InsertClient(req.Name, req.CNPJ)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

// UpdateClientHandler updates client fields.
func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := UpdateClientDetails(id, req.Name, req.CNPJ)
	if err != nil {
		apierr.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(c)
}

// DeleteClientHandler handles client removal.
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := RemoveClient(id)
	if err != nil {
		apierr.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
