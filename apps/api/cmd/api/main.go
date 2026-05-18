package main

import (
	"log"

	"github.com/meetext/backend/internal/app"
	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	l := logger.New(cfg.Log.Level, cfg.Log.Pretty, cfg.Log.File, cfg.Log.MaxSizeMB)

	l.Info().
		Str("app", cfg.App.Name).
		Str("env", cfg.App.Env).
		Str("http_addr", cfg.HTTP.Host+":"+cfg.HTTP.Port).
		Str("frontend_url", cfg.App.FrontendURL).
		Bool("google_client_id_set", cfg.OAuth.GoogleClientID != "").
		Bool("google_client_secret_set", cfg.OAuth.GoogleClientSecret != "").
		Str("google_redirect_url", cfg.OAuth.GoogleRedirectURL).
		Msg("configuration loaded")

	if err := cfg.ValidateAPI(); err != nil {
		l.Fatal().Err(err).Msg("invalid API configuration")
	}

	a, err := app.New(cfg, l)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to initialize app")
	}

	if err := a.Run(); err != nil {
		l.Fatal().Err(err).Msg("app exited with error")
	}
}
