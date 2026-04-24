package http

import (
	"encoding/json"
	"net/http"

	"github.com/task-manager/api/pkg/validator"
)

const maxJSONBodyBytes = 1 << 20

// ReadAndValidateJSON decodes JSON from r.Body (size-limited) into v and runs validator.Validate.
// On failure it writes the appropriate 400 response and returns false.
func ReadAndValidateJSON(w http.ResponseWriter, r *http.Request, v any) bool {
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		Error(w, http.StatusBadRequest, "invalid request body")
		return false
	}
	if err := validator.Validate(v); err != nil {
		Error(w, http.StatusBadRequest, "validation failed: "+err.Error())
		return false
	}
	return true
}
