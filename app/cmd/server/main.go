package main

import (
	"fmt"
	"os"

	"github.com/task-manager/api/cmd/server/bootstrap"
	"github.com/task-manager/api/config"
	"github.com/task-manager/api/pkg/logger"
)

func main() {
	cfg := config.Load()
	if cfg.DatabaseURL == "" {
		fmt.Fprintln(os.Stderr, "DATABASE_URL is required")
		os.Exit(1)
	}
	if cfg.JWTSecret == "" {
		fmt.Fprintln(os.Stderr, "JWT_SECRET is required")
		os.Exit(1)
	}

	log := logger.New(logger.Config{
		Level: cfg.LogLevel,
		JSON:  cfg.LogJSON,
	})

	if err := bootstrap.Run(cfg, log); err != nil {
		log.Fatal().Err(err).Msg("application failed")
	}
}
