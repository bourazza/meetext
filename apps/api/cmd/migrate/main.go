package main

import (
	"flag"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/meetext/backend/internal/config"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up | down")
	flag.Parse()

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	migrationsPath := "file://migrations"
	if p := os.Getenv("MIGRATIONS_PATH"); p != "" {
		migrationsPath = "file://" + p
	}

	m, err := migrate.New(migrationsPath, cfg.DB.DSN)
	if err != nil {
		log.Fatalf("migrate: init: %v", err)
	}
	defer m.Close()

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate up: %v", err)
		}
		log.Println("migrations applied")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migrate down: %v", err)
		}
		log.Println("migrations rolled back")
	default:
		log.Fatalf("unknown direction: %s", *direction)
	}
}
