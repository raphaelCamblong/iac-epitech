package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"github.com/task-manager/api/pkg/logger"
)

const (
	HeaderCorrelationID = "correlation_id"
	HeaderTraceID       = "trace_id"
)

// Logging returns a middleware that logs requests with correlation_id and trace_id.
func Logging(log zerolog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			correlationID := r.Header.Get(HeaderCorrelationID)
			if correlationID == "" {
				correlationID = uuid.New().String()
			}
			traceID := r.Header.Get(HeaderTraceID)
			if traceID == "" {
				traceID = uuid.New().String()
			}
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()
			reqID := middleware.GetReqID(r.Context())
			if reqID == "" {
				reqID = traceID
			}
			reqLog := logger.WithContext(log, correlationID, traceID)
			ctx := reqLog.WithContext(r.Context())
			r = r.WithContext(ctx)
			ww.Header().Set(HeaderCorrelationID, correlationID)
			next.ServeHTTP(ww, r)
			reqLog.Info().
				Str("request_id", reqID).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int("status", ww.Status()).
				Int("bytes", ww.BytesWritten()).
				Dur("duration_ms", time.Since(start)).
				Msg("request completed")
		})
	}
}
