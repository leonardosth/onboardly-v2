package user

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/go-chi/chi/v5"
	"onboardly-backend/internal/apierr"
	"onboardly-backend/internal/auth"
)

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// validatePassword enforces at least 8 chars, 1 letter, and 1 number.
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

// ListUsersHandler handles GET /api/users
func ListUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := ListUsers()
	if err != nil {
		apierr.WriteError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	if users == nil {
		users = []auth.User{} // return empty array instead of null
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// CreateUserHandler handles POST /api/users
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apierr.WriteError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if !validatePassword(req.Password) {
		apierr.WriteError(w, http.StatusBadRequest, "A senha deve conter no mínimo 8 caracteres, 1 letra e 1 número")
		return
	}

	user, err := auth.CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		// Basic check for duplicate email (assumes unique constraint error)
		if err.Error() == "invalid role, must be Admin or Analista" {
			apierr.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		apierr.WriteError(w, http.StatusConflict, "Email já cadastrado ou erro ao criar usuário")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// DeleteUserHandler handles DELETE /api/users/{id}
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	// Try to get current user ID to prevent self-deletion
	// Note: Our JWT currently doesn't store user ID, only email and role.
	// We will look up the current user's ID by email, or skip this if not easily available.
	// For now, pass empty string since we rely on the client to block self deletion or we extract email.
	
	currentEmail := r.Context().Value(auth.UserEmailKey).(string)
	
	// Quick hack to get current user ID by email to prevent self-deletion
	// We can use ListUsers and find the one matching currentEmail
	var currentUserID string
	users, _ := ListUsers()
	for _, u := range users {
		if u.Email == currentEmail {
			currentUserID = u.ID
			break
		}
	}

	err := DeleteUser(id, currentUserID)
	if err != nil {
		if err.Error() == "cannot delete your own account" || err.Error() == "cannot delete the last admin in the system" {
			apierr.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err.Error() == "user not found" {
			apierr.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		apierr.WriteError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
