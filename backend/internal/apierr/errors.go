package apierr

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse defines the standard JSON payload structure for API errors.
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteError sends a formatted JSON error response to the client.
func WriteError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}
