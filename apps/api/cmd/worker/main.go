package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	l := logger.New(cfg.Log.Level, cfg.Log.Pretty, cfg.Log.File, cfg.Log.MaxSizeMB)
	l.Info().Msg("worker starting")

	// TODO: initialize Redis queue client and register job handlers
	// e.g. transcription_worker, extraction_worker, export_worker

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info().Msg("worker stopped")
}
