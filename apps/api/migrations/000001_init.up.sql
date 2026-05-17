-- =========================================
-- MEETEXT DATABASE SCHEMA (MVP VERSION)
-- PostgreSQL
-- =========================================

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- =========================================
-- ENUMS
-- =========================================

CREATE TYPE subscription_plan AS ENUM ('free', 'pro', 'business');
CREATE TYPE workspace_role    AS ENUM ('owner', 'admin', 'member');
CREATE TYPE upload_type       AS ENUM ('audio', 'video', 'pdf', 'docx');
CREATE TYPE meeting_status    AS ENUM ('uploaded', 'processing', 'completed', 'failed', 'needs_review');
CREATE TYPE task_status       AS ENUM ('todo', 'in_progress', 'review', 'done');
CREATE TYPE task_priority     AS ENUM ('low', 'medium', 'high', 'urgent');
CREATE TYPE project_status    AS ENUM ('planning', 'active', 'review', 'completed');
CREATE TYPE document_type     AS ENUM ('summary', 'requirements', 'technical_doc', 'sprint_plan', 'client_notes', 'decision_log');

-- =========================================
-- USERS
-- =========================================

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name     VARCHAR(255) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    avatar_url    TEXT,
    plan          subscription_plan DEFAULT 'free',
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- WORKSPACES
-- =========================================

CREATE TABLE workspaces (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    owner_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    name       VARCHAR(255) NOT NULL,
    logo_url   TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- WORKSPACE MEMBERS
-- =========================================

CREATE TABLE workspace_members (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    user_id      UUID REFERENCES users(id) ON DELETE CASCADE,
    role         workspace_role DEFAULT 'member',
    created_at   TIMESTAMP DEFAULT NOW(),
    UNIQUE(workspace_id, user_id)
);

-- =========================================
-- CLIENTS
-- =========================================

CREATE TABLE clients (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id  UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    company_name  VARCHAR(255) NOT NULL,
    contact_name  VARCHAR(255),
    contact_email VARCHAR(255),
    logo_url      TEXT,
    notes         TEXT,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- PROJECTS
-- =========================================

CREATE TABLE projects (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    client_id    UUID REFERENCES clients(id) ON DELETE SET NULL,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    status       project_status DEFAULT 'planning',
    progress     INTEGER DEFAULT 0,
    start_date   DATE,
    end_date     DATE,
    created_at   TIMESTAMP DEFAULT NOW(),
    updated_at   TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- MEETINGS
-- =========================================

CREATE TABLE meetings (
    id                      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id            UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id              UUID REFERENCES projects(id) ON DELETE SET NULL,
    client_id               UUID REFERENCES clients(id) ON DELETE SET NULL,
    title                   VARCHAR(255) NOT NULL,
    upload_type             upload_type NOT NULL,
    original_file_url       TEXT NOT NULL,
    transcript              TEXT,
    ai_summary              TEXT,
    duration_seconds        INTEGER,
    language                VARCHAR(50),
    status                  meeting_status DEFAULT 'uploaded',
    processing_started_at   TIMESTAMP,
    processing_completed_at TIMESTAMP,
    uploaded_by             UUID REFERENCES users(id),
    created_at              TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- MEETING PARTICIPANTS
-- =========================================

CREATE TABLE meeting_participants (
    id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id        UUID REFERENCES meetings(id) ON DELETE CASCADE,
    participant_name  VARCHAR(255) NOT NULL,
    participant_email VARCHAR(255),
    created_at        TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- TASKS
-- =========================================

CREATE TABLE tasks (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id  UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id    UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id    UUID REFERENCES meetings(id) ON DELETE SET NULL,
    title         VARCHAR(255) NOT NULL,
    description   TEXT,
    status        task_status DEFAULT 'todo',
    priority      task_priority DEFAULT 'medium',
    due_date      DATE,
    ai_generated  BOOLEAN DEFAULT TRUE,
    ai_confidence DECIMAL(5,2),
    assigned_to   UUID REFERENCES users(id),
    created_at    TIMESTAMP DEFAULT NOW(),
    updated_at    TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- GOALS
-- =========================================

CREATE TABLE goals (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id   UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id   UUID REFERENCES meetings(id) ON DELETE SET NULL,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    completed    BOOLEAN DEFAULT FALSE,
    target_date  DATE,
    created_at   TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- DEADLINES
-- =========================================

CREATE TABLE deadlines (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id   UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id   UUID REFERENCES meetings(id) ON DELETE SET NULL,
    title        VARCHAR(255) NOT NULL,
    description  TEXT,
    due_date     DATE NOT NULL,
    completed    BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- DECISIONS
-- =========================================

CREATE TABLE decisions (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id  UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id    UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id    UUID REFERENCES meetings(id) ON DELETE CASCADE,
    decision_text TEXT NOT NULL,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- RISKS / BLOCKERS
-- =========================================

CREATE TABLE blockers (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id   UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id   UUID REFERENCES meetings(id) ON DELETE CASCADE,
    blocker_text TEXT NOT NULL,
    severity     VARCHAR(50),
    resolved     BOOLEAN DEFAULT FALSE,
    created_at   TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- DOCUMENTS
-- =========================================

CREATE TABLE documents (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id    UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    project_id      UUID REFERENCES projects(id) ON DELETE CASCADE,
    meeting_id      UUID REFERENCES meetings(id) ON DELETE SET NULL,
    title           VARCHAR(255) NOT NULL,
    type            document_type NOT NULL,
    content         TEXT,
    generated_by_ai BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT NOW(),
    updated_at      TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- INTEGRATIONS
-- =========================================

CREATE TABLE integrations (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id  UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    provider      VARCHAR(100) NOT NULL,
    access_token  TEXT,
    refresh_token TEXT,
    is_active     BOOLEAN DEFAULT TRUE,
    created_at    TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- EXPORTS
-- =========================================

CREATE TABLE exports (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    workspace_id UUID REFERENCES workspaces(id) ON DELETE CASCADE,
    meeting_id   UUID REFERENCES meetings(id) ON DELETE CASCADE,
    export_type  VARCHAR(50),
    file_url     TEXT,
    created_by   UUID REFERENCES users(id),
    created_at   TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- AI PROCESSING LOGS
-- =========================================

CREATE TABLE ai_processing_logs (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    meeting_id UUID REFERENCES meetings(id) ON DELETE CASCADE,
    step_name  VARCHAR(255),
    status     VARCHAR(50),
    message    TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- NOTIFICATIONS
-- =========================================

CREATE TABLE notifications (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID REFERENCES users(id) ON DELETE CASCADE,
    title      VARCHAR(255),
    body       TEXT,
    is_read    BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- =========================================
-- INDEXES
-- =========================================

CREATE INDEX idx_meetings_workspace  ON meetings(workspace_id);
CREATE INDEX idx_projects_workspace  ON projects(workspace_id);
CREATE INDEX idx_tasks_project       ON tasks(project_id);
CREATE INDEX idx_documents_project   ON documents(project_id);
CREATE INDEX idx_meetings_status     ON meetings(status);
CREATE INDEX idx_tasks_status        ON tasks(status);
CREATE INDEX idx_notifications_user  ON notifications(user_id);
