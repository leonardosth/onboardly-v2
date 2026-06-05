package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"onboardly-backend/internal/auth"
	"onboardly-backend/internal/client"
	"onboardly-backend/internal/config"
	"onboardly-backend/internal/dashboard"
	"onboardly-backend/internal/meeting"
	"onboardly-backend/internal/project"
	"onboardly-backend/internal/user"
)

// NewRouter initializes the chi router, mounts standard middlewares and API handlers.
func NewRouter(cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(corsMiddleware)

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// API routes
	r.Route("/api", func(r chi.Router) {
		// Public Auth routes
		r.Post("/auth/register", auth.RegisterHandler)
		r.Post("/auth/login", auth.LoginHandler(cfg))

		// Webhook client sync endpoint (secured internally via token check)
		r.Post("/webhooks/clients", client.WebhookSyncHandler(cfg))

		// Authenticated routes group
		r.Group(func(r chi.Router) {
			r.Use(auth.AuthMiddleware(cfg.JWTSecret))

			// Dashboard analytics
			r.Get("/dashboard", dashboard.DashboardHandler)

			// Client routes
			r.Get("/clients", client.ListClientsHandler)
			r.Get("/clients/{id}", client.GetClientHandler)
			r.Post("/clients", client.CreateClientHandler)
			r.Put("/clients/{id}", client.UpdateClientHandler)

			// Admin-only deletion endpoint
			r.With(auth.RequireRole("Admin")).Delete("/clients/{id}", client.DeleteClientHandler)

			// Project routes
			r.Get("/projects", project.ListProjectsHandler)
			r.Get("/projects/{id}", project.GetProjectHandler)
			r.Post("/projects", project.CreateProjectHandler)
			r.Put("/projects/{id}", project.UpdateProjectStatusHandler)

			// Meeting routes
			r.Get("/meetings", meeting.ListMeetingsHandler)
			r.Post("/meetings", meeting.CreateMeetingHandler)

			// User routes (Admin only)
			r.Group(func(r chi.Router) {
				r.Use(auth.RequireRole("Admin"))
				r.Get("/users", user.ListUsersHandler)
				r.Post("/users", user.CreateUserHandler)
				r.Delete("/users/{id}", user.DeleteUserHandler)
			})
		})
	})

	return r
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Webhook-Token")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
