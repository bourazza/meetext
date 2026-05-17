# Meetext — Backend API

> Go REST API for the Meetext AI meeting intelligence platform.

---

## Table of Contents

1. [Tech Stack](#tech-stack)
2. [Project Structure](#project-structure)
3. [Architecture](#architecture)
4. [Getting Started](#getting-started)
   - [Local (no Docker)](#local-no-docker)
   - [Docker](#docker)
5. [Environment Variables](#environment-variables)
6. [Database](#database)
7. [API Reference](#api-reference)
   - [Health](#health)
   - [Auth](#auth)
   - [Workspaces](#workspaces)
   - [Meetings](#meetings)
8. [Response Format](#response-format)
9. [Error Codes](#error-codes)
10. [Middleware](#middleware)
11. [Domain Model](#domain-model)
12. [Makefile Commands](#makefile-commands)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.22 |
| Router | Chi v5 |
| Database | PostgreSQL 16 |
| DB Driver | pgx v5 |
| Config | Viper |
| Logging | Zerolog |
| Validation | go-playground/validator v10 |
| Auth | JWT (golang-jwt/jwt v5) |
| Password | bcrypt |
| Migrations | golang-migrate v4 |
| Storage | Local filesystem (S3-ready interface) |
| Containerization | Docker + Docker Compose |

---

## Project Structure

```
apps/api/
├── cmd/
│   ├── api/main.go           # API server entrypoint
│   ├── migrate/main.go       # Database migration runner
│   └── worker/main.go        # Background worker entrypoint (stub)
│
├── internal/
│   ├── app/
│   │   └── app.go            # Dependency wiring + HTTP server lifecycle
│   │
│   ├── config/
│   │   ├── config.go         # Strongly typed config structs
│   │   └── env.go            # Viper loader with defaults
│   │
│   ├── domain/               # Pure business entities + repository interfaces
│   │   ├── user/
│   │   ├── workspace/
│   │   ├── meeting/
│   │   ├── client/
│   │   ├── project/
│   │   ├── task/
│   │   ├── document/
│   │   ├── goal/
│   │   ├── deadline/
│   │   ├── decision/
│   │   └── blocker/
│   │
│   ├── usecase/              # Application business logic
│   │   ├── auth/             # Register, login, refresh token
│   │   ├── workspace/        # Workspace + member management
│   │   └── meeting/          # File upload, list, get, delete
│   │
│   ├── repository/
│   │   └── postgres/         # pgx implementations of domain interfaces
│   │       ├── user_repo.go
│   │       ├── workspace_repo.go
│   │       └── meeting_repo.go
│   │
│   ├── infrastructure/
│   │   ├── auth/jwt.go       # JWT issue + validate
│   │   ├── db/postgres.go    # pgx connection pool
│   │   ├── password/         # bcrypt hash + compare
│   │   └── storage/          # File storage (local provider)
│   │
│   └── delivery/http/
│       ├── handler/          # HTTP request handlers
│       ├── middleware/        # Auth + logger middleware
│       └── router/           # Chi router + route registration
│
├── migrations/
│   ├── 000001_init.up.sql    # Full schema (15 tables, 8 enums)
│   └── 000001_init.down.sql  # Full rollback
│
├── pkg/                      # Shared utilities (no business logic)
│   ├── apperr/               # Custom error types + HTTP status mapping
│   ├── response/             # Standard JSON response envelope
│   ├── logger/               # Zerolog factory
│   ├── validator/            # JSON decode + struct validation
│   ├── utils/                # Random hex, generic pointer helper
│   └── constants/            # Context keys, roles, MIME types, limits
│
├── .env.example              # All environment variables documented
├── .air.toml                 # Hot reload config (Air)
├── Dockerfile                # Multi-stage build
├── entrypoint.sh             # Runs migrations then starts server
├── go.mod
└── sqlc.yaml                 # sqlc code generation config
```

---

## Architecture

The backend follows **Clean Architecture** with strict layer separation:

```
Delivery (HTTP handlers)
    ↓  calls
Use Cases (business logic)
    ↓  calls
Domain interfaces (repository contracts)
    ↓  implemented by
Infrastructure (postgres, storage, jwt...)
```

**Rules:**
- Domain layer has zero external dependencies
- Use cases depend only on domain interfaces, never on concrete implementations
- Handlers never touch the database directly
- All dependencies are injected via constructors — no globals

---

## Getting Started

### Local (no Docker)

**Prerequisites:** Go 1.22+, PostgreSQL running locally

**1. Create the database**
```bash
make local-setup
```
This creates the `meetext` postgres user and `meetext` database.

If it fails on password, run manually:
```bash
sudo -u postgres psql -c "ALTER USER meetext WITH PASSWORD 'meetext';"
```

**2. Copy and configure environment**
```bash
cp apps/api/.env.example apps/api/.env
```

**3. Run migrations**
```bash
make local-migrate
```

**4. Start the server**
```bash
make local-run
```

Server starts at `http://localhost:8080`

**Or do all steps at once:**
```bash
make local-dev
```

---

### Docker

**Prerequisites:** Docker + Docker Compose

```bash
make docker-up
```

This will:
- Build the Go binary inside a multi-stage Docker image
- Start PostgreSQL (mapped to port **5433** to avoid conflict with local postgres)
- Start Redis (port 6379)
- Run database migrations automatically via `entrypoint.sh`
- Start the API server on port **8080**

```bash
make docker-logs    # tail live logs
make docker-down    # stop all containers
make docker-reset   # wipe volumes + restart fresh
```

> **Note:** Docker postgres is exposed on port `5433` externally. The API container connects internally on `5432` via the Docker network.

---

## Environment Variables

All variables are in `apps/api/.env.example`:

### App
| Variable | Default | Description |
|---|---|---|
| `APP_NAME` | `meetext` | Application name |
| `APP_ENV` | `development` | Environment: `development` / `production` |
| `APP_VERSION` | `0.1.0` | API version |
| `FRONTEND_URL` | `http://localhost:3000` | Allowed CORS origin |

### HTTP Server
| Variable | Default | Description |
|---|---|---|
| `HTTP_HOST` | `0.0.0.0` | Bind address |
| `HTTP_PORT` | `8080` | Listen port |
| `HTTP_READ_TIMEOUT` | `15s` | Request read timeout |
| `HTTP_WRITE_TIMEOUT` | `15s` | Response write timeout |
| `HTTP_IDLE_TIMEOUT` | `60s` | Keep-alive idle timeout |

### Database
| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | `postgres://meetext:meetext@localhost:5432/meetext?sslmode=disable` | Full PostgreSQL DSN |
| `DB_MAX_OPEN_CONNS` | `25` | Max open connections in pool |
| `DB_MAX_IDLE_CONNS` | `5` | Min idle connections in pool |
| `DB_MAX_LIFETIME` | `5m` | Max connection lifetime |

### JWT
| Variable | Default | Description |
|---|---|---|
| `JWT_ACCESS_SECRET` | — | **Required.** Min 32 chars. Signs access tokens |
| `JWT_REFRESH_SECRET` | — | **Required.** Min 32 chars. Signs refresh tokens |
| `JWT_ACCESS_TTL` | `15m` | Access token lifetime |
| `JWT_REFRESH_TTL` | `168h` | Refresh token lifetime (7 days) |

### Storage
| Variable | Default | Description |
|---|---|---|
| `STORAGE_PROVIDER` | `local` | `local` / `s3` / `supabase` |
| `STORAGE_LOCAL_PATH` | `./uploads` | Local upload directory |
| `STORAGE_BUCKET` | — | S3 bucket name |
| `STORAGE_REGION` | — | S3 region |
| `STORAGE_ACCESS_KEY` | — | S3 access key |
| `STORAGE_SECRET_KEY` | — | S3 secret key |
| `STORAGE_ENDPOINT` | — | Custom S3 endpoint (Supabase/R2) |

### Redis
| Variable | Default | Description |
|---|---|---|
| `REDIS_ADDR` | `localhost:6379` | Redis address |
| `REDIS_PASSWORD` | — | Redis password (optional) |
| `REDIS_DB` | `0` | Redis database index |

### AI
| Variable | Default | Description |
|---|---|---|
| `OLLAMA_URL` | `http://localhost:11434` | Ollama LLM server URL |
| `OLLAMA_MODEL` | `llama3` | Model to use for extraction |
| `WHISPER_URL` | `http://localhost:9000` | Whisper transcription server URL |

### Logging
| Variable | Default | Description |
|---|---|---|
| `LOG_LEVEL` | `info` | `debug` / `info` / `warn` / `error` |
| `LOG_PRETTY` | `true` | `true` for colored console, `false` for JSON |

---

## Database

### Schema Overview

The database has **15 tables** and **8 enums**:

**Enums**
| Enum | Values |
|---|---|
| `subscription_plan` | `free`, `pro`, `business` |
| `workspace_role` | `owner`, `admin`, `member` |
| `upload_type` | `audio`, `video`, `pdf`, `docx` |
| `meeting_status` | `uploaded`, `processing`, `completed`, `failed`, `needs_review` |
| `task_status` | `todo`, `in_progress`, `review`, `done` |
| `task_priority` | `low`, `medium`, `high`, `urgent` |
| `project_status` | `planning`, `active`, `review`, `completed` |
| `document_type` | `summary`, `requirements`, `technical_doc`, `sprint_plan`, `client_notes`, `decision_log` |

**Tables**
| Table | Description |
|---|---|
| `users` | Platform accounts with subscription plan |
| `workspaces` | Multi-tenant workspace owned by a user |
| `workspace_members` | User ↔ workspace membership with role |
| `clients` | External clients linked to a workspace |
| `projects` | Projects inside a workspace, optionally linked to a client |
| `meetings` | Uploaded meeting files (audio/video/pdf/docx) |
| `meeting_participants` | Named participants extracted from a meeting |
| `tasks` | AI-generated or manual tasks linked to a project/meeting |
| `goals` | Goals extracted from meetings |
| `deadlines` | Deadlines extracted from meetings |
| `decisions` | Decisions recorded from meetings |
| `blockers` | Risks and blockers identified in meetings |
| `documents` | AI-generated documents (summaries, specs, etc.) |
| `integrations` | OAuth tokens for Notion, Jira, etc. |
| `exports` | Export history (PDF, DOCX, Sheets) |
| `ai_processing_logs` | Step-by-step AI pipeline audit log |
| `notifications` | In-app notifications per user |

### Migrations

```bash
make migrate-up          # apply all pending migrations
make migrate-down        # roll back all migrations
make migrate-create NAME=add_something   # create new migration file
```

Migration files live in `apps/api/migrations/` and follow the `golang-migrate` naming convention:
```
000001_init.up.sql
000001_init.down.sql
```

---

## API Reference

**Base URL:** `http://localhost:8080`

**All protected routes require:**
```
Authorization: Bearer <access_token>
```

---

### Health

#### `GET /health`

Check if the server is running.

**Auth:** None

**Response `200`**
```json
{
  "success": true,
  "data": { "status": "ok" }
}
```

---

### Auth

#### `POST /api/v1/auth/register`

Register a new user. Automatically creates a workspace and adds the user as owner.

**Auth:** None

**Request Body**
```json
{
  "full_name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword",
  "workspace_name": "My Agency"
}
```

| Field | Type | Required | Rules |
|---|---|---|---|
| `full_name` | string | ✓ | min 2, max 100 |
| `email` | string | ✓ | valid email |
| `password` | string | ✓ | min 8 chars |
| `workspace_name` | string | ✓ | min 2, max 100 |

**Response `201`**
```json
{
  "success": true,
  "data": {
    "user": {
      "id": "uuid",
      "full_name": "John Doe",
      "email": "john@example.com",
      "plan": "free",
      "created_at": "2026-01-01T00:00:00Z"
    },
    "workspace": {
      "id": "uuid",
      "name": "My Agency",
      "owner_id": "uuid",
      "created_at": "2026-01-01T00:00:00Z"
    },
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `CONFLICT` | 409 | Email already registered |
| `VALIDATION_ERROR` | 422 | Invalid input fields |

---

#### `POST /api/v1/auth/login`

Authenticate an existing user.

**Auth:** None

**Request Body**
```json
{
  "email": "john@example.com",
  "password": "securepassword"
}
```

**Response `200`**
```json
{
  "success": true,
  "data": {
    "user": { ... },
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `INVALID_CREDENTIALS` | 401 | Wrong email or password |
| `VALIDATION_ERROR` | 422 | Missing fields |

---

#### `POST /api/v1/auth/refresh`

Get a new access token using a valid refresh token.

**Auth:** None

**Request Body**
```json
{
  "refresh_token": "eyJ..."
}
```

**Response `200`**
```json
{
  "success": true,
  "data": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `TOKEN_EXPIRED` | 401 | Refresh token has expired |
| `TOKEN_INVALID` | 401 | Refresh token is malformed |

---

### Workspaces

All workspace routes require a valid `Authorization: Bearer <token>` header.

---

#### `GET /api/v1/workspaces`

List all workspaces the authenticated user belongs to.

**Response `200`**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "owner_id": "uuid",
      "name": "My Agency",
      "logo_url": null,
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

#### `GET /api/v1/workspaces/{workspaceID}`

Get a single workspace by ID.

**Path Params**
| Param | Type | Description |
|---|---|---|
| `workspaceID` | UUID | Workspace ID |

**Response `200`**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "owner_id": "uuid",
    "name": "My Agency",
    "logo_url": null,
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

---

#### `PATCH /api/v1/workspaces/{workspaceID}`

Update workspace name. Requires `owner` or `admin` role.

**Path Params**
| Param | Type | Description |
|---|---|---|
| `workspaceID` | UUID | Workspace ID |

**Request Body**
```json
{
  "name": "New Workspace Name"
}
```

**Response `200`** — updated workspace object

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `FORBIDDEN` | 403 | Requester is not owner or admin |
| `NOT_FOUND` | 404 | Workspace not found |

---

#### `GET /api/v1/workspaces/{workspaceID}/members`

List all members of a workspace.

**Response `200`**
```json
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "workspace_id": "uuid",
      "user_id": "uuid",
      "role": "owner",
      "created_at": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

#### `DELETE /api/v1/workspaces/{workspaceID}/members/{userID}`

Remove a member from a workspace. Requires `owner` or `admin` role.

**Path Params**
| Param | Type | Description |
|---|---|---|
| `workspaceID` | UUID | Workspace ID |
| `userID` | UUID | User to remove |

**Response `204`** — no content

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `FORBIDDEN` | 403 | Requester is not owner or admin |
| `NOT_FOUND` | 404 | Member not found |

---

### Meetings

All meeting routes require a valid `Authorization: Bearer <token>` header.

---

#### `POST /api/v1/workspaces/{workspaceID}/meetings`

Upload a meeting file. Accepts `multipart/form-data`.

**Path Params**
| Param | Type | Description |
|---|---|---|
| `workspaceID` | UUID | Workspace ID |

**Form Fields**
| Field | Type | Required | Description |
|---|---|---|---|
| `file` | file | ✓ | The meeting file |
| `title` | string | — | Defaults to filename if omitted |
| `project_id` | UUID | — | Link to an existing project |
| `client_id` | UUID | — | Link to an existing client |

**Supported MIME types**
| MIME | Upload Type |
|---|---|
| `audio/mpeg` | audio |
| `audio/wav` | audio |
| `video/mp4` | video |
| `application/pdf` | pdf |

**Max file size:** 500 MB

**Response `201`**
```json
{
  "success": true,
  "data": {
    "id": "uuid",
    "workspace_id": "uuid",
    "project_id": null,
    "client_id": null,
    "title": "Q1 Planning Call",
    "upload_type": "audio",
    "original_file_url": "workspaces/.../meetings/....mp3",
    "transcript": null,
    "ai_summary": null,
    "duration_seconds": null,
    "language": null,
    "status": "uploaded",
    "uploaded_by": "uuid",
    "created_at": "2026-01-01T00:00:00Z"
  }
}
```

**Errors**
| Code | Status | Meaning |
|---|---|---|
| `FILE_TOO_LARGE` | 413 | File exceeds 500 MB |
| `UNSUPPORTED_FILE` | 400 | MIME type not allowed |
| `MISSING_FILE` | 400 | No file field in form |

---

#### `GET /api/v1/workspaces/{workspaceID}/meetings`

List all meetings in a workspace. Supports pagination.

**Query Params**
| Param | Type | Default | Description |
|---|---|---|---|
| `limit` | int | `20` | Max results (capped at 100) |
| `offset` | int | `0` | Pagination offset |

**Response `200`**
```json
{
  "success": true,
  "data": [ { ...meeting }, { ...meeting } ]
}
```

---

#### `GET /api/v1/workspaces/{workspaceID}/meetings/{meetingID}`

Get a single meeting by ID.

**Path Params**
| Param | Type | Description |
|---|---|---|
| `workspaceID` | UUID | Workspace ID |
| `meetingID` | UUID | Meeting ID |

**Response `200`** — full meeting object including `transcript` and `ai_summary` when available.

---

#### `DELETE /api/v1/workspaces/{workspaceID}/meetings/{meetingID}`

Delete a meeting and its uploaded file.

**Response `204`** — no content

---

## Response Format

Every response follows the same envelope:

**Success**
```json
{
  "success": true,
  "data": { ... }
}
```

**Error**
```json
{
  "success": false,
  "error": {
    "code": "NOT_FOUND",
    "message": "resource not found"
  }
}
```

**Validation Error**
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "validation failed",
    "fields": {
      "email": "must be a valid email address",
      "password": "value is too short (min 8)"
    }
  }
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|---|---|---|
| `NOT_FOUND` | 404 | Resource does not exist |
| `UNAUTHORIZED` | 401 | Missing or invalid auth header |
| `FORBIDDEN` | 403 | Authenticated but not permitted |
| `CONFLICT` | 409 | Resource already exists (e.g. duplicate email) |
| `BAD_REQUEST` | 400 | Malformed request or invalid UUID |
| `VALIDATION_ERROR` | 422 | Input failed validation rules |
| `INVALID_CREDENTIALS` | 401 | Wrong email or password |
| `TOKEN_EXPIRED` | 401 | JWT access/refresh token has expired |
| `TOKEN_INVALID` | 401 | JWT token is malformed or tampered |
| `FILE_TOO_LARGE` | 413 | Upload exceeds 500 MB limit |
| `UNSUPPORTED_FILE` | 400 | MIME type not accepted |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

---

## Middleware

All requests pass through the following middleware stack (in order):

| Middleware | Description |
|---|---|
| `RequestID` | Attaches a unique `X-Request-ID` to every request |
| `RealIP` | Extracts real client IP from `X-Forwarded-For` |
| `Logger` | Logs method, path, status, latency, request ID via Zerolog |
| `Recoverer` | Catches panics and returns `500` instead of crashing |
| `RateLimit` | 100 requests per minute per IP |
| `CORS` | Allows requests only from `FRONTEND_URL` |
| `Auth` | Validates `Bearer` JWT on protected routes only |

---

## Domain Model

```
User
 └── belongs to many Workspaces (via workspace_members)

Workspace
 ├── has many Members (workspace_members)
 ├── has many Clients
 ├── has many Projects
 └── has many Meetings

Project
 ├── belongs to Workspace
 ├── optionally belongs to Client
 ├── has many Meetings
 ├── has many Tasks
 ├── has many Goals
 ├── has many Deadlines
 ├── has many Decisions
 ├── has many Blockers
 └── has many Documents

Meeting
 ├── belongs to Workspace + Project + Client
 ├── has inline Transcript (text)
 ├── has inline AI Summary (text)
 ├── has many Participants
 ├── has many Tasks (ai_generated)
 ├── has many Goals
 ├── has many Decisions
 ├── has many Blockers
 └── has many Documents
```

---

## Makefile Commands

### Local Development
```bash
make local-setup      # create postgres user + database
make local-migrate    # run migrations against local postgres
make local-run        # start Go server locally
make local-dev        # setup + migrate + run (full bootstrap)
make local-reset      # drop DB, recreate, re-migrate
```

### Docker
```bash
make docker-up        # build + start all services
make docker-down      # stop all services
make docker-logs      # tail API container logs
make docker-reset     # wipe volumes + restart fresh
make docker-db-only   # start only postgres + redis (run Go locally)
```

### Build
```bash
make build            # compile production binary to bin/api
```

### Testing & Quality
```bash
make test             # run all tests with race detector
make test-cover       # run tests + open HTML coverage report
make lint             # run golangci-lint
make fmt              # gofmt all files
make vet              # go vet
make tidy             # go mod tidy + verify
```

### Database
```bash
make migrate-up                        # apply pending migrations
make migrate-down                      # roll back all migrations
make migrate-create NAME=add_projects  # create new migration file
```

### sqlc
```bash
make sqlc             # generate query code from SQL
make sqlc-verify      # verify queries compile
```
