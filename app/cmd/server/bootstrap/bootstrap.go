package bootstrap

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"

	"github.com/task-manager/api/config"
	authctrl "github.com/task-manager/api/internal/adapter/http/auth"
	"github.com/task-manager/api/internal/adapter/http/middleware"
	taskctrl "github.com/task-manager/api/internal/adapter/http/task"
	"github.com/task-manager/api/internal/adapter/repository"
	authuc "github.com/task-manager/api/internal/usecase/auth"
	taskuc "github.com/task-manager/api/internal/usecase/task"
	"github.com/task-manager/api/pkg/jwtauth"
)

// Run starts the application.
func Run(cfg *config.Config, log zerolog.Logger) error {
	db, err := repository.NewGormDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	sqlDB, _ := db.DB()
	defer sqlDB.Close()

	taskRepo := repository.NewTaskGormRepository(db)
	userRepo := repository.NewUserGormRepository(db)

	taskUC := taskuc.NewTaskUseCase(taskRepo)
	authUC := authuc.NewAuthUseCase(userRepo)

	signer := jwtauth.NewSigner(cfg.JWTSecret, cfg.JWTExpiry)
	validator := jwtauth.NewValidator(cfg.JWTSecret)

	authHandler := authctrl.NewHandler(authUC, signer)
	taskHandler := taskctrl.NewHandler(taskUC)

	router := NewRouter(RouterConfig{
		AuthHandler: authHandler,
		TaskHandler: taskHandler,
		Auth:        middleware.Auth(validator, log),
		Recovery:    middleware.Recovery(log),
		Logging:     middleware.Logging(log),
		RateLimit:   middleware.RateLimit(cfg.RateLimitRPS, cfg.RateLimitBurst),
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}
	go func() {
		log.Info().Str("port", cfg.Port).Msg("server starting")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down")
	ctx := context.Background()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("shutdown error")
	}
	log.Info().Msg("server stopped")
	return nil
}
