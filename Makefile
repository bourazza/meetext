API_DIR := ./apps/api
BINARY  := bin/api

# Load local .env for DB vars
-include $(API_DIR)/.env
export

.PHONY: help \
        local-setup local-db local-migrate local-run local-reset \
        docker-up docker-down docker-logs docker-reset docker-migrate \
        build test test-cover lint tidy fmt vet \
        migrate-up migrate-down migrate-create sqlc

# ─────────────────────────────────────────────────────────────────────────────
help: ## Show all available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2}'

# ─────────────────────────────────────────────────────────────────────────────
# LOCAL (no Docker)
# ─────────────────────────────────────────────────────────────────────────────

local-setup: ## Create local postgres user + database
	@echo "→ Creating postgres user 'meetext'..."
	sudo -u postgres psql -c "CREATE USER meetext WITH PASSWORD 'meetext';" 2>/dev/null || true
	@echo "→ Creating database 'meetext'..."
	sudo -u postgres createdb -O meetext meetext 2>/dev/null || true
	@echo "✓ Local database ready"

local-migrate: ## Run migrations against local postgres
	cd $(API_DIR) && go run ./cmd/migrate -direction=up

local-run: ## Start API server locally (no Docker)
	cd $(API_DIR) && go run ./cmd/api

local-reset: ## Drop and recreate local database, re-run migrations
	@echo "→ Dropping database 'meetext'..."
	sudo -u postgres dropdb --if-exists meetext
	@echo "→ Recreating database 'meetext'..."
	sudo -u postgres createdb -O meetext meetext
	@echo "→ Running migrations..."
	cd $(API_DIR) && go run ./cmd/migrate -direction=up
	@echo "✓ Local database reset complete"

local-dev: ## Setup DB + migrate + run server (full local bootstrap)
	$(MAKE) local-setup
	$(MAKE) local-migrate
	$(MAKE) local-run

# ─────────────────────────────────────────────────────────────────────────────
# DOCKER
# ─────────────────────────────────────────────────────────────────────────────

docker-up: ## Start all services with Docker (postgres, redis, api)
	docker compose up -d --build

docker-down: ## Stop all Docker services
	docker compose down

docker-logs: ## Tail API logs from Docker
	docker compose logs -f api

docker-migrate: ## Run migrations against Docker postgres (port 5433)
	DATABASE_URL=postgres://meetext:meetext@localhost:5433/meetext?sslmode=disable \
		cd $(API_DIR) && go run ./cmd/migrate -direction=up

docker-reset: ## Destroy Docker volumes and restart fresh
	docker compose down -v
	docker compose up -d --build

docker-db-only: ## Start only postgres + redis (useful for local Go dev with Docker DB)
	docker compose up -d postgres redis

# ─────────────────────────────────────────────────────────────────────────────
# BUILD
# ─────────────────────────────────────────────────────────────────────────────

build: ## Build production binary to bin/api
	@mkdir -p bin
	cd $(API_DIR) && CGO_ENABLED=0 go build -ldflags="-s -w" -o ../../$(BINARY) ./cmd/api
	@echo "✓ Binary built at $(BINARY)"

# ─────────────────────────────────────────────────────────────────────────────
# TESTING & QUALITY
# ─────────────────────────────────────────────────────────────────────────────

test: ## Run all tests
	cd $(API_DIR) && go test ./... -v -race -timeout 60s

test-cover: ## Run tests with HTML coverage report
	cd $(API_DIR) && go test ./... -coverprofile=coverage.out
	cd $(API_DIR) && go tool cover -html=coverage.out

lint: ## Run golangci-lint
	cd $(API_DIR) && golangci-lint run ./...

tidy: ## Tidy go modules
	cd $(API_DIR) && go mod tidy && go mod verify

fmt: ## Format all Go files
	cd $(API_DIR) && gofmt -w .

vet: ## Run go vet
	cd $(API_DIR) && go vet ./...

# ─────────────────────────────────────────────────────────────────────────────
# MIGRATIONS
# ─────────────────────────────────────────────────────────────────────────────

migrate-up: ## Apply all pending migrations
	cd $(API_DIR) && go run ./cmd/migrate -direction=up

migrate-down: ## Roll back all migrations
	cd $(API_DIR) && go run ./cmd/migrate -direction=down

migrate-create: ## Create new migration file (usage: make migrate-create NAME=add_users)
	migrate create -ext sql -dir $(API_DIR)/migrations -seq $(NAME)

# ─────────────────────────────────────────────────────────────────────────────
# SQLC
# ─────────────────────────────────────────────────────────────────────────────

sqlc: ## Generate sqlc query code
	cd $(API_DIR) && sqlc generate

sqlc-verify: ## Verify sqlc queries
	cd $(API_DIR) && sqlc verify
