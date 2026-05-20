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
	"github.com/meetext/backend/internal/infrastructure/email"
	infraoauth "github.com/meetext/backend/internal/infrastructure/oauth"
	infraollama "github.com/meetext/backend/internal/infrastructure/ollama"
	infrapdf "github.com/meetext/backend/internal/infrastructure/pdf"
	"github.com/meetext/backend/internal/infrastructure/storage"
	infrawhisper "github.com/meetext/backend/internal/infrastructure/whisper"
	"github.com/meetext/backend/internal/repository/postgres"
	ucai "github.com/meetext/backend/internal/usecase/ai"
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
	authTokenRepo := postgres.NewAuthTokenRepository(pool)
	workspaceRepo := postgres.NewWorkspaceRepository(pool)
	meetingRepo := postgres.NewMeetingRepository(pool)

	// Use cases
	emailSvc := email.NewLogService(log)
	authUC := ucauth.NewUseCase(userRepo, workspaceRepo, authTokenRepo, jwtSvc, emailSvc, cfg.App.FrontendURL, cfg.Auth.RequireEmailVerified)
	workspaceUC := ucworkspace.NewUseCase(workspaceRepo)

	// OAuth providers
	googleProvider := infraoauth.NewGoogle(cfg.OAuth)
	githubProvider := infraoauth.NewGitHub(cfg.OAuth)

	if !googleProvider.IsConfigured() {
		log.Warn().Msg("app: Google OAuth disabled; GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET, or GOOGLE_REDIRECT_URL is not set")
	} else {
		log.Info().
			Str("provider", string(googleProvider.Name())).
			Str("redirect_url", cfg.OAuth.GoogleRedirectURL).
			Msg("app: OAuth provider registered")
	}
	if !githubProvider.IsConfigured() {
		log.Warn().Msg("app: GitHub OAuth disabled; GITHUB_CLIENT_ID, GITHUB_CLIENT_SECRET, or GITHUB_REDIRECT_URL is not set")
	} else {
		log.Info().
			Str("provider", string(githubProvider.Name())).
			Str("redirect_url", cfg.OAuth.GitHubRedirectURL).
			Msg("app: OAuth provider registered")
	}

	// AI Providers
	ollamaProvider := infraollama.NewProvider(cfg.AI, log)
	whisperProvider := infrawhisper.NewMockProvider(log)
	aiUC := ucai.NewUseCase(ollamaProvider, log)
	pdfExtractor := infrapdf.NewExtractor(log)
	meetingUC := ucmeeting.NewUseCase(meetingRepo, store, aiUC, pdfExtractor)

	// Handlers
	handlers := router.Handlers{
		Auth:        handler.NewAuthHandler(authUC, log),
		OAuthGoogle: handler.NewOAuthHandler(googleProvider, authUC, cfg.App.FrontendURL, log),
		OAuthGitHub: handler.NewOAuthHandler(githubProvider, authUC, cfg.App.FrontendURL, log),
		Workspace:   handler.NewWorkspaceHandler(workspaceUC),
		Meeting:     handler.NewMeetingHandler(meetingUC, log),
		AI:          handler.NewAIHandler(aiUC, whisperProvider, log),
	}

	httpHandler := router.New(log, jwtSvc, authTokenRepo, cfg.App.FrontendURL, handlers)

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
