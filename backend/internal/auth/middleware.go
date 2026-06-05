package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"onboardly-backend/internal/apierr"
)

type contextKey string

// Context keys for injecting auth data into request context.
const (
	UserEmailKey contextKey = "user_email"
	UserRoleKey  contextKey = "user_role"
)

// AuthMiddleware extracts, validates, and injects JWT claims into the request context.
func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				apierr.WriteError(w, http.StatusUnauthorized, "authorization header missing")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				apierr.WriteError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			tokenStr := parts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				apierr.WriteError(w, http.StatusUnauthorized, "invalid token")
				return
			}

			ctx := context.WithValue(r.Context(), UserEmailKey, claims.Email)
			ctx = context.WithValue(ctx, UserRoleKey, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole restricts access to HTTP routes to specific predefined user roles.
func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roleVal := r.Context().Value(UserRoleKey)
			if roleVal == nil {
				apierr.WriteError(w, http.StatusForbidden, "access denied: role missing")
				return
			}

			userRole := roleVal.(string)
			allowed := false
			for _, r := range roles {
				if userRole == r {
					allowed = true
					break
				}
			}

			if !allowed {
				apierr.WriteError(w, http.StatusForbidden, "access denied: insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
