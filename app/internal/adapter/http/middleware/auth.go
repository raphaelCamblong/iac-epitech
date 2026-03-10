package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/rs/zerolog"

	"github.com/task-manager/api/pkg/jwtauth"
)

type contextKey string

const userContextKey contextKey = "user"

// UserInfo holds authenticated user data from JWT.
type UserInfo struct {
	UserID string
	Email  string
}

// Auth validates JWT Bearer token and injects UserInfo into context.
// Fail-closed: returns 401 on missing or invalid auth.
func Auth(validator *jwtauth.Validator, log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				log.Debug().Str("path", r.URL.Path).Msg("missing authorization header")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(auth, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				log.Debug().Msg("invalid authorization format")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			token := parts[1]
			claims, err := validator.Validate(token)
			if err != nil {
				log.Debug().Err(err).Msg("invalid token")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), userContextKey, &UserInfo{
				UserID: claims.UserID,
				Email:  claims.Email,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserFromContext returns the authenticated user from context.
func UserFromContext(ctx context.Context) *UserInfo {
	u, _ := ctx.Value(userContextKey).(*UserInfo)
	return u
}

// UserContextKey returns the context key for UserInfo. Used by tests to inject user.
func UserContextKey() interface{} {
	return userContextKey
}
