package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/internal/delivery/http/handler"
	"github.com/meetext/backend/internal/delivery/http/router"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/internal/infrastructure/db"
	"github.com/meetext/backend/internal/infrastructure/storage"
	"github.com/meetext/backend/internal/repository/postgres"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	ucmeeting "github.com/meetext/backend/internal/usecase/meeting"
	ucworkspace "github.com/meetext/backend/internal/usecase/workspace"
	"github.com/rs/zerolog"
)

type App struct {
	cfg    *config.Config
	log    zerolog.Logger
	server *http.Server
}

func New(cfg *config.Config, log zerolog.Logger) (*App, error) {
	ctx := context.Background()

	// Infrastructure
	pool, err := db.NewPool(ctx, cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("app: db pool: %w", err)
	}

	jwtSvc := infraauth.NewJWTService(cfg.JWT)

	store, err := storage.NewLocalProvider(cfg.Storage.LocalPath, cfg.App.FrontendURL+"/uploads")
	if err != nil {
		return nil, fmt.Errorf("app: storage: %w", err)
	}

	// Repositories
	userRepo := postgres.NewUserRepository(pool)
	workspaceRepo := postgres.NewWorkspaceRepository(pool)
	meetingRepo := postgres.NewMeetingRepository(pool)

	// Use cases
	authUC := ucauth.NewUseCase(userRepo, workspaceRepo, jwtSvc)
	workspaceUC := ucworkspace.NewUseCase(workspaceRepo)
	meetingUC := ucmeeting.NewUseCase(meetingRepo, store)

	// Handlers
	handlers := router.Handlers{
		Auth:      handler.NewAuthHandler(authUC),
		Workspace: handler.NewWorkspaceHandler(workspaceUC),
		Meeting:   handler.NewMeetingHandler(meetingUC),
	}

	httpHandler := router.New(log, jwtSvc, cfg.App.FrontendURL, handlers)

	srv := &http.Server{
		Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler:      httpHandler,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
		IdleTimeout:  cfg.HTTP.IdleTimeout,
	}

	return &App{cfg: cfg, log: log, server: srv}, nil
}

func (a *App) Run() error {
	errCh := make(chan error, 1)

	go func() {
		a.log.Info().Str("addr", a.server.Addr).Msg("server starting")
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return fmt.Errorf("server error: %w", err)
	case sig := <-quit:
		a.log.Info().Str("signal", sig.String()).Msg("shutting down")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("graceful shutdown failed: %w", err)
	}

	a.log.Info().Msg("server stopped")
	return nil
}
