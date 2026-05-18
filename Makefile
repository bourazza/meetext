API_DIR := ./apps/api
WEB_DIR := ./apps/web

BINARY := ./bin/api
PID_DIR := ./.pids

API_PORT := 8080
WEB_PORT := 3000

# Load env file if exists
ifneq (,$(wildcard $(API_DIR)/.env))
	include $(API_DIR)/.env
	export
endif

SHELL := /bin/bash

.PHONY: help start stop \
	local-setup local-migrate local-run local-reset local-dev \
	web-install web-dev web-build web-start web-lint \
	docker-up docker-down docker-logs docker-reset docker-migrate docker-db-only \
	build test test-cover lint tidy fmt vet \
	migrate-up migrate-down migrate-create sqlc sqlc-verify clean

# ─────────────────────────────────────────────────────────────
# HELP
# ─────────────────────────────────────────────────────────────

help: ## Show commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-24s\033[0m %s\n", $$1, $$2}'

# ─────────────────────────────────────────────────────────────
# INTERNAL HELPERS
# ─────────────────────────────────────────────────────────────

define kill_port
	@if lsof -ti:$(1) >/dev/null 2>&1; then \
		lsof -ti:$(1) | xargs kill -9; \
	fi
endef

define kill_pid_file
	@if [ -f $(1) ]; then \
		pkill -P $$(cat $(1)) 2>/dev/null || true; \
		kill $$(cat $(1)) 2>/dev/null || true; \
		rm -f $(1); \
	fi
endef

# ─────────────────────────────────────────────────────────────
# START / STOP
# ─────────────────────────────────────────────────────────────

start: ## Start API + Web locally
	@mkdir -p $(PID_DIR)

	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "Starting Meetext local environment"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

	$(call kill_pid_file,$(PID_DIR)/api.pid)
	$(call kill_pid_file,$(PID_DIR)/web.pid)

	$(call kill_port,$(API_PORT))

	@for port in 3000 3001 3002 3003 3004 3005; do \
		lsof -ti:$$port 2>/dev/null | xargs kill -9 2>/dev/null || true; \
	done
	@sleep 2

	@echo "→ Running migrations..."
	@cd $(API_DIR) && go run ./cmd/migrate -direction=up

	@echo "→ Building API..."
	@mkdir -p ./bin
	@cd $(API_DIR) && go build -o ../../$(BINARY) ./cmd/api

	@echo "→ Starting API..."
	@sh scripts/start-api.sh

	@echo "→ Waiting for API..."
	@for i in $$(seq 1 20); do \
		curl -sf http://localhost:$(API_PORT)/health >/dev/null 2>&1 && break || sleep 1; \
	done

	@echo "✓ API running"

	@echo "→ Starting Web..."
	@rm -rf $(WEB_DIR)/.next
	@sh scripts/start-web.sh

	@echo "→ Waiting for Web..."
	@for i in $$(seq 1 90); do \
		nc -z localhost $(WEB_PORT) 2>/dev/null && break; \
		sleep 1; \
	done
	@kill -0 $$(cat $(PID_DIR)/web.pid) 2>/dev/null || (echo "ERROR: Web process died. Check /tmp/meetext-web.log"; cat /tmp/meetext-web.log | tail -20; exit 1)

	@echo ""
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
	@echo "✓ API : http://localhost:$(API_PORT)"
	@echo "✓ Web : http://localhost:$(WEB_PORT)"
	@echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

stop: ## Stop all local services
	@echo "→ Stopping services..."
	@if [ -f $(PID_DIR)/api.pid ]; then kill -9 $$(cat $(PID_DIR)/api.pid) 2>/dev/null || true; rm -f $(PID_DIR)/api.pid; fi
	@if [ -f $(PID_DIR)/web.pid ]; then pkill -P $$(cat $(PID_DIR)/web.pid) 2>/dev/null || true; kill -9 $$(cat $(PID_DIR)/web.pid) 2>/dev/null || true; rm -f $(PID_DIR)/web.pid; fi
	@for port in 8080 3000 3001 3002 3003 3004 3005; do lsof -ti:$$port 2>/dev/null | xargs kill -9 2>/dev/null || true; done
	@echo "✓ All services stopped"
# ─────────────────────────────────────────────────────────────
# LOCAL
# ─────────────────────────────────────────────────────────────

local-setup: ## Create local postgres DB
	@echo "→ Creating postgres user..."
	@sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname='meetext'" | grep -q 1 || \
	sudo -u postgres psql -c "CREATE USER meetext WITH PASSWORD 'meetext';"

	@echo "→ Creating database..."
	@sudo -u postgres psql -lqt | cut -d \| -f 1 | grep -qw meetext || \
	sudo -u postgres createdb -O meetext meetext

	@echo "✓ Database ready"

local-migrate: ## Run migrations
	cd $(API_DIR) && go run ./cmd/migrate -direction=up

local-run: ## Run API locally
	cd $(API_DIR) && go run ./cmd/api

local-reset: ## Reset database
	@sudo -u postgres dropdb --if-exists meetext
	@sudo -u postgres createdb -O meetext meetext
	@cd $(API_DIR) && go run ./cmd/migrate -direction=up
	@echo "✓ Database reset complete"

local-dev: ## Full local bootstrap
	@$(MAKE) local-setup
	@$(MAKE) local-migrate
	@$(MAKE) local-run

# ─────────────────────────────────────────────────────────────
# WEB
# ─────────────────────────────────────────────────────────────

web-install:
	cd $(WEB_DIR) && npm install

web-dev:
	cd $(WEB_DIR) && npm run dev

web-build:
	cd $(WEB_DIR) && npm run build

web-start:
	cd $(WEB_DIR) && npm run start

web-lint:
	cd $(WEB_DIR) && npm run lint

# ─────────────────────────────────────────────────────────────
# DOCKER
# ─────────────────────────────────────────────────────────────

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f api

docker-migrate:
	cd $(API_DIR) && \
	DATABASE_URL="postgres://meetext:meetext@localhost:5433/meetext?sslmode=disable" \
	go run ./cmd/migrate -direction=up

docker-reset:
	docker compose down -v
	docker compose up -d --build

docker-db-only:
	docker compose up -d postgres redis

# ─────────────────────────────────────────────────────────────
# BUILD
# ─────────────────────────────────────────────────────────────

build: ## Build production binary
	@mkdir -p ./bin
	@cd $(API_DIR) && \
	CGO_ENABLED=0 go build -ldflags="-s -w" -o ../../$(BINARY) ./cmd/api

	@echo "✓ Binary built: $(BINARY)"

clean: ## Remove generated files
	rm -rf bin
	rm -rf $(PID_DIR)

# ─────────────────────────────────────────────────────────────
# TESTING
# ─────────────────────────────────────────────────────────────

test:
	cd $(API_DIR) && go test ./... -v -race -timeout=60s

test-cover:
	cd $(API_DIR) && go test ./... -coverprofile=coverage.out
	cd $(API_DIR) && go tool cover -html=coverage.out

lint:
	cd $(API_DIR) && golangci-lint run ./...

tidy:
	cd $(API_DIR) && go mod tidy && go mod verify

fmt:
	cd $(API_DIR) && gofmt -w .

vet:
	cd $(API_DIR) && go vet ./...

# ─────────────────────────────────────────────────────────────
# MIGRATIONS
# ─────────────────────────────────────────────────────────────

migrate-up:
	cd $(API_DIR) && go run ./cmd/migrate -direction=up

migrate-down:
	cd $(API_DIR) && go run ./cmd/migrate -direction=down

migrate-create:
	migrate create -ext sql -dir $(API_DIR)/migrations -seq $(NAME)

# ─────────────────────────────────────────────────────────────
# SQLC
# ─────────────────────────────────────────────────────────────

sqlc:
	cd $(API_DIR) && sqlc generate

sqlc-verify:
	cd $(API_DIR) && sqlc verify