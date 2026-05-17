package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/meetext/backend/internal/delivery/http/handler"
	httpmiddleware "github.com/meetext/backend/internal/delivery/http/middleware"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/pkg/response"
	"github.com/rs/zerolog"
	"time"
)

type Handlers struct {
	Auth      *handler.AuthHandler
	Workspace *handler.WorkspaceHandler
	Meeting   *handler.MeetingHandler
}

func New(log zerolog.Logger, jwt *infraauth.JWTService, frontendURL string, h Handlers) http.Handler {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(httpmiddleware.Logger(log))
	r.Use(chimiddleware.Recoverer)
	r.Use(httprate.LimitByIP(100, time.Minute))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		// Public auth routes
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", h.Auth.Register)
			r.Post("/login", h.Auth.Login)
			r.Post("/refresh", h.Auth.Refresh)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.Auth(jwt))

			r.Route("/workspaces", func(r chi.Router) {
				r.Get("/", h.Workspace.List)
				r.Get("/{workspaceID}", h.Workspace.GetByID)
				r.Patch("/{workspaceID}", h.Workspace.Update)
				r.Get("/{workspaceID}/members", h.Workspace.ListMembers)
				r.Delete("/{workspaceID}/members/{userID}", h.Workspace.RemoveMember)

				r.Route("/{workspaceID}/meetings", func(r chi.Router) {
					r.Post("/", h.Meeting.Upload)
					r.Get("/", h.Meeting.List)
					r.Get("/{meetingID}", h.Meeting.GetByID)
					r.Delete("/{meetingID}", h.Meeting.Delete)
				})
			})
		})
	})

	return r
}
