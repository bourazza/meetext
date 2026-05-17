package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/lumberjack.v2"
)

// New creates a zerolog.Logger that writes to:
//   - stdout (pretty colored output in development, JSON in production)
//   - a rotating log file at cfg.File
func New(level string, pretty bool, logFile string, maxSizeMB int) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}

	zerolog.SetGlobalLevel(lvl)
	zerolog.TimeFieldFormat = time.RFC3339

	// Ensure log directory exists
	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "logger: failed to create log dir: %v\n", err)
		}
	}

	fileWriter := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    maxSizeMB, // MB before rotation
		MaxBackups: 7,         // keep last 7 rotated files
		MaxAge:     30,        // days
		Compress:   true,      // gzip rotated files
	}

	// File always gets JSON (machine-readable)
	jsonWriter := zerolog.New(fileWriter).With().Timestamp().Logger()
	_ = jsonWriter

	var consoleWriter io.Writer
	if pretty {
		consoleWriter = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	} else {
		consoleWriter = os.Stdout
	}

	// Multi-writer: console + file
	multi := zerolog.MultiLevelWriter(consoleWriter, fileWriter)

	logger := zerolog.New(multi).With().Timestamp().Logger()

	// Replace global logger too so any log.Info() calls also go to file
	log.Logger = logger

	return logger
}
