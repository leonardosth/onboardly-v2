package auth

import (
	"encoding/json"
	"net/http"

	"onboardly-backend/internal/apierr"
	"onboardly-backend/internal/config"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

// RegisterHandler handles user signup requests.
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Role == "" {
		apierr.WriteError(w, http.StatusBadRequest, "email, password, and role are required")
		return
	}

	user, err := CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)
}

// LoginHandler handles user authentication requests.
func LoginHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apierr.WriteError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		if req.Email == "" || req.Password == "" {
			apierr.WriteError(w, http.StatusBadRequest, "email and password are required")
			return
		}

		token, role, err := AuthenticateUser(req.Email, req.Password, cfg.JWTSecret)
		if err != nil {
			apierr.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(loginResponse{
			Token: token,
			Role:  role,
		})
	}
}
