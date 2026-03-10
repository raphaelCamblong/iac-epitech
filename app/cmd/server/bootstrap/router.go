package bootstrap

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	authctrl "github.com/task-manager/api/internal/adapter/http/auth"
	taskctrl "github.com/task-manager/api/internal/adapter/http/task"
)

// RouterConfig holds router dependencies.
type RouterConfig struct {
	AuthHandler *authctrl.Handler
	TaskHandler *taskctrl.Handler
	Auth        func(http.Handler) http.Handler
	Recovery    func(http.Handler) http.Handler
	Logging     func(http.Handler) http.Handler
	RateLimit   func(http.Handler) http.Handler
}

// NewRouter builds the HTTP router.
func NewRouter(cfg RouterConfig) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(cfg.Recovery)
	r.Use(cfg.Logging)
	r.Use(cfg.RateLimit)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", cfg.AuthHandler.Register)
		r.Post("/login", cfg.AuthHandler.Login)
	})

	r.Route("/tasks", func(r chi.Router) {
		r.Use(cfg.Auth)
		r.Post("/", cfg.TaskHandler.Create)
		r.Get("/", cfg.TaskHandler.List)
		r.Get("/{id}", cfg.TaskHandler.Get)
		r.Put("/{id}", cfg.TaskHandler.Update)
		r.Delete("/{id}", cfg.TaskHandler.Delete)
	})

	return r
}
