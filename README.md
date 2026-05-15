# Meetext

> AI-powered meeting intelligence and project documentation platform.

Meetext transforms meetings, recordings, audio files, videos, and PDFs into structured project intelligence using AI.

Generated outputs include:

* AI summaries
* Tasks and tickets
* Goals
* Deadlines
* Decisions
* Risks and blockers
* Technical documentation
* Exportable project artifacts

---

# Full System Architecture Documentation

## Overview

Meetext is an AI-powered meeting intelligence and project documentation platform.

The system transforms:

* Audio meetings
* Video meetings
* PDFs
* Client recordings
* Voice notes

into:

* Structured documentation
* Tasks/tickets
* Goals
* Deadlines
* Decisions
* Risks/blockers
* Exportable project artifacts

The platform is designed around:

* Clean Architecture
* Domain-driven modularity
* AI orchestration
* Background processing
* Scalability
* Low operational cost during MVP stage

---

# Table of Contents

1. High-Level Architecture
2. Core Technology Stack
3. Product Modules
4. AI Pipeline
5. Why AI Logic Should NOT Be Inside n8n
6. n8n Responsibilities
7. Backend Clean Architecture
8. Queue System
9. Database Design Philosophy
10. Storage Architecture
11. Security Architecture
12. Frontend Architecture
13. API Design
14. Export System
15. Suggested MVP Scope
16. Scaling Strategy
17. Recommended AI Models
18. DevOps Architecture
19. Observability
20. Final Architectural Philosophy

---

# 1. High-Level Architecture

```text
User
  ↓
Next.js Frontend
  ↓
Go API Gateway
  ↓
---------------------------------
| Core Backend Services        |
|                              |
| - Auth                       |
| - Meetings                   |
| - Projects                   |
| - Tasks                      |
| - Documents                  |
| - AI Orchestration           |
| - File Processing            |
---------------------------------
  ↓
Job Queue (Redis)
  ↓
Workers
  ↓
---------------------------------
| AI Layer                     |
|                              |
| - Whisper                    |
| - Ollama                     |
| - Prompt Engine              |
---------------------------------
  ↓
PostgreSQL Database
  ↓
Storage (S3/Supabase)
  ↓
n8n Integrations Layer
```

---

# 2. Core Technology Stack

## Frontend

| Technology    | Purpose                 |
| ------------- | ----------------------- |
| Next.js       | Main frontend framework |
| React         | UI rendering            |
| Tailwind CSS  | Styling                 |
| Zustand/Redux | State management        |
| React Query   | API caching             |
| shadcn/ui     | UI components           |

---

## Backend

| Technology     | Purpose               |
| -------------- | --------------------- |
| Go             | Main backend language |
| Gin/Echo/Fiber | HTTP framework        |
| JWT            | Authentication        |
| Redis          | Queue and caching     |
| PostgreSQL     | Main database         |
| Docker         | Containerization      |

---

## AI Stack

| Technology         | Purpose              |
| ------------------ | -------------------- |
| Whisper            | Speech-to-text       |
| Ollama             | Local LLM runtime    |
| Llama/Qwen/Mistral | AI extraction models |

---

## Automation

| Technology | Purpose                      |
| ---------- | ---------------------------- |
| n8n        | Integrations and automations |

---

# 3. Product Modules

## Authentication Module

Responsibilities:

* User registration
* Login/logout
* JWT generation
* Workspace permissions
* Session validation

Endpoints:

* POST /auth/register
* POST /auth/login
* POST /auth/refresh
* POST /auth/logout

Security:

* Password hashing using bcrypt
* JWT access tokens
* Refresh tokens
* Middleware authorization

---

## Workspace Module

Responsibilities:

* Multi-tenant workspace isolation
* Team management
* Subscription plan management

Key Concepts:

* Every resource belongs to a workspace
* Users can belong to multiple workspaces
* Roles:

  * owner
  * admin
  * member

---

## Meetings Module

Responsibilities:

* Upload meetings
* Store recordings
* Trigger AI processing
* Manage transcripts
* Store AI summaries

Supported Upload Types:

* mp3
* wav
* mp4
* pdf
* docx

Flow:

```text
Upload File
   ↓
Validate File
   ↓
Store File
   ↓
Create Meeting Record
   ↓
Push AI Job to Queue
```

---

## AI Processing Module

This is the core intelligence layer of Meetext.

Responsibilities:

* Speech transcription
* Transcript cleaning
* Structured extraction
* JSON normalization
* Summary generation
* Entity extraction

---

# 4. AI Pipeline

## Step 1 — Upload

Frontend uploads file.

Backend:

* validates type
* validates size
* stores file
* creates DB record
* pushes processing job

---

## Step 2 — Audio Extraction

If video:

```text
Video → Extract Audio → WAV
```

Tools:

* ffmpeg

---

## Step 3 — Whisper Transcription

Worker sends audio to Whisper.

Whisper returns:

* transcript
* timestamps
* language detection

Stored in:

* meetings.transcript

---

## Step 4 — Transcript Cleaning

Cleanup layer:

* remove noise
* normalize punctuation
* merge broken sentences
* remove filler words

Optional future:

* speaker diarization

---

## Step 5 — LLM Extraction

Transcript sent to Ollama.

Prompt examples:

* tasks extraction
* goals extraction
* deadlines extraction
* blockers extraction
* decisions extraction

Expected Output:

```json
{
  "tasks": [],
  "goals": [],
  "deadlines": [],
  "decisions": [],
  "blockers": []
}
```

---

## Step 6 — Validation Layer

AI output validated.

Rules:

* valid JSON only
* remove hallucinations
* deduplicate tasks
* normalize dates
* normalize priorities

---

## Step 7 — Persistence

Structured entities inserted into database.

Tables:

* tasks
* goals
* deadlines
* blockers
* decisions
* documents

---

# 5. Why AI Logic Should NOT Be Inside n8n

AI extraction belongs in backend services because:

| Reason            | Explanation                   |
| ----------------- | ----------------------------- |
| Scalability       | Easier worker scaling         |
| Reliability       | Better retries and validation |
| Performance       | Better memory handling        |
| Maintainability   | Easier testing                |
| Security          | Sensitive logic isolated      |
| Prompt management | Centralized prompts           |

n8n should only handle:

* integrations
* exports
* notifications
* workflow automations

---

# 6. n8n Responsibilities

## Recommended n8n Workflows

### Export Workflow

Triggered when:

* user requests export

Actions:

* generate PDF
* generate DOCX
* upload export
* notify user

---

### Notion Sync Workflow

Triggered when:

* project updated
* document generated

Actions:

* create/update Notion pages

---

### Jira Sync Workflow

Triggered when:

* tasks generated

Actions:

* create Jira tickets

---

### Notification Workflow

Actions:

* send emails
* send Slack notifications
* webhook callbacks

---

# 7. Backend Clean Architecture

## Architecture Layers

```text
Interfaces Layer
    ↓
Use Cases Layer
    ↓
Domain Layer
    ↓
Infrastructure Layer
```

---

## Domain Layer

Contains:

* business entities
* interfaces
* business rules

Example:

```text
meeting/
 ├── entity.go
 ├── repository.go
 └── service.go
```

Responsibilities:

* no external dependencies
* pure business logic

---

## Use Case Layer

Contains application workflows.

Examples:

* upload meeting
* process transcript
* generate tasks
* export project

This layer orchestrates:

* repositories
* AI services
* queues
* validation

---

## Infrastructure Layer

Contains external implementations.

Examples:

* PostgreSQL
* Redis
* Whisper
* Ollama
* S3
* JWT

---

## Interfaces Layer

Handles:

* HTTP requests
* DTO mapping
* validation
* responses

Contains:

* handlers
* middleware
* routes

---

# 8. Queue System

Queue system responsibilities:

* background jobs
* retries
* async processing
* heavy AI workloads

Recommended:

* Redis queues

Workers:

| Worker               | Responsibility      |
| -------------------- | ------------------- |
| transcription_worker | Whisper jobs        |
| extraction_worker    | LLM extraction      |
| export_worker        | PDF/DOCX generation |

---

# 9. Database Design Philosophy

Core Principles:

* Multi-tenant architecture
* Workspace isolation
* AI entity normalization
* Scalable relational design

Main Entities:

* users
* workspaces
* meetings
* projects
* tasks
* documents
* goals
* deadlines
* blockers
* decisions

---

# 10. Storage Architecture

## File Types

Stored Files:

* original uploads
* extracted audio
* exports
* generated docs

---

## Storage Providers

MVP:

* Supabase Storage

Scalable:

* AWS S3
* Cloudflare R2

---

# 11. Security Architecture

## Authentication

* JWT access token
* Refresh token
* HTTP-only cookies

---

## Authorization

Workspace-based RBAC.

Roles:

* owner
* admin
* member

---

## Upload Security

Validate:

* MIME types
* file size
* malicious uploads

---

## AI Security

Prevent:

* prompt injection
* malformed JSON
* AI hallucinations

---

# 12. Frontend Architecture

## Routing

Using Next.js App Router.

Pages:

* dashboard
* meetings
* projects
* tasks
* documents
* clients
* settings

---

## State Management

Recommended:

* Zustand for global state
* React Query for server state

---

## UI Design Principles

Meetext UI should be:

* minimal
* calm
* productivity-focused
* AI-native
* not overloaded

Avoid:

* enterprise complexity
* unnecessary analytics
* chatbot-heavy interfaces

---

# 13. API Design

## REST API Structure

```text
/api/v1/auth
/api/v1/meetings
/api/v1/projects
/api/v1/tasks
/api/v1/documents
/api/v1/clients
```

---

## Response Structure

```json
{
  "success": true,
  "data": {},
  "message": "success"
}
```

---

## Error Structure

```json
{
  "success": false,
  "error": {
    "code": "INVALID_FILE",
    "message": "Unsupported file type"
  }
}
```

---

# 14. Export System

Supported Exports:

* PDF
* DOCX
* Markdown
* JSON
* Google Docs
* Google Sheets
* Notion

Premium Exports:

* Jira
* Linear
* Notion sync

---

# 15. Suggested MVP Scope

## MVP Features

Required:

* Upload meetings
* Whisper transcription
* AI summaries
* Task extraction
* Goals extraction
* Documents page
* PDF export
* Google Sheets export

Do NOT build initially:

* AI chat
* real-time collaboration
* browser extension
* mobile app
* advanced analytics

---

# 16. Scaling Strategy

## Stage 1 — Local Development

Environment:

* CPU only
* local Ollama
* local PostgreSQL

Goal:

* validate workflows

---

## Stage 2 — Small VPS

Environment:

* single server
* Docker compose
* small Redis instance

Goal:

* first real users

---

## Stage 3 — Production Scale

Environment:

* dedicated workers
* GPU inference
* distributed queues
* Kubernetes

Goal:

* scale AI workloads

---

# 17. Recommended AI Models

## Transcription

| Model         | Use                |
| ------------- | ------------------ |
| Whisper Tiny  | Fast testing       |
| Whisper Base  | MVP                |
| Whisper Large | Production quality |

---

## Extraction Models

| Model   | Use                          |
| ------- | ---------------------------- |
| Qwen    | Strong structured extraction |
| Llama 3 | General reasoning            |
| Mistral | Lightweight extraction       |

---

# 18. DevOps Architecture

## Docker Services

Recommended services:

```text
- web
- api
- worker
- postgres
- redis
- ollama
- nginx
- n8n
```

---

# 19. Observability

Recommended:

* structured logging
* request tracing
* queue monitoring
* AI processing metrics

Tools:

* Grafana
* Prometheus
* Loki

---

# 20. Final Architectural Philosophy

Meetext should be built as:

"An AI-powered project documentation operating system."

Core principles:

* AI-first
* workflow simplicity
* structured outputs
* automation-friendly
* scalable backend
* lightweight UX
* modular architecture

The most important long-term asset is NOT transcription.

The real value is:

* structured project intelligence
* workflow automation
* persistent organizational memory
* turning meetings into actionable systems
