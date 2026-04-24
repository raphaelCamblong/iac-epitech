package http

import (
	"errors"
	"net/http"

	"github.com/task-manager/api/internal/port"
)

// WritePortError maps known port layer errors to HTTP responses.
// Returns true if err was recognized and a response was written.
func WritePortError(w http.ResponseWriter, err error) bool {
	switch {
	case errors.Is(err, port.ErrNotFound):
		Error(w, http.StatusNotFound, "task not found")
		return true
	case errors.Is(err, port.ErrConflict):
		Error(w, http.StatusConflict, "timestamp or concurrency conflict")
		return true
	case errors.Is(err, port.ErrDuplicateUser):
		Error(w, http.StatusConflict, "user already exists")
		return true
	case errors.Is(err, port.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, "invalid credentials")
		return true
	default:
		return false
	}
}
