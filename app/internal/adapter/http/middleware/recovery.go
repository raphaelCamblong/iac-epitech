package middleware

import (
	"log/slog"
	"net/http"

	"github.com/rs/zerolog"
)

// Recovery recovers from panics and returns 500. Never let panic propagate.
func Recovery(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error().Interface("panic", err).Msg("panic recovered")
					slog.Default().Error("panic recovered", "err", err)
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
