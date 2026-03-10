package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config holds logger configuration.
type Config struct {
	Level string
	JSON  bool
}

// New creates a zerolog.Logger with the given config.
func New(cfg Config) zerolog.Logger {
	var out io.Writer = os.Stdout
	if !cfg.JSON {
		out = zerolog.ConsoleWriter{Out: os.Stdout}
	}
	level := zerolog.InfoLevel
	if l, err := zerolog.ParseLevel(cfg.Level); err == nil {
		level = l
	}
	return zerolog.New(out).Level(level).With().Timestamp().Logger()
}

// WithContext returns a logger with correlation_id and trace_id.
func WithContext(l zerolog.Logger, correlationID, traceID string) zerolog.Logger {
	ctx := l.With()
	if correlationID != "" {
		ctx = ctx.Str("correlation_id", correlationID)
	}
	if traceID != "" {
		ctx = ctx.Str("trace_id", traceID)
	}
	return ctx.Logger()
}

// L returns the default global logger.
func L() *zerolog.Logger {
	return &log.Logger
}
