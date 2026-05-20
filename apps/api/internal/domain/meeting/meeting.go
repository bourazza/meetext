package meeting

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type UploadType string

const (
	UploadTypeAudio UploadType = "audio"
	UploadTypeVideo UploadType = "video"
	UploadTypePDF   UploadType = "pdf"
	UploadTypeDOCX  UploadType = "docx"
)

type Status string

const (
	StatusUploaded    Status = "uploaded"
	StatusProcessing  Status = "processing"
	StatusCompleted   Status = "completed"
	StatusFailed      Status = "failed"
	StatusNeedsReview Status = "needs_review"
)

type Meeting struct {
	ID                     uuid.UUID  `json:"id"`
	WorkspaceID            uuid.UUID  `json:"workspace_id"`
	ProjectID              *uuid.UUID `json:"project_id,omitempty"`
	ClientID               *uuid.UUID `json:"client_id,omitempty"`
	Title                  string     `json:"title"`
	UploadType             UploadType `json:"upload_type"`
	OriginalFileURL        string     `json:"original_file_url"`
	Transcript             *string    `json:"transcript,omitempty"`
	AISummary              *string    `json:"ai_summary,omitempty"`
	DurationSeconds        *int       `json:"duration_seconds,omitempty"`
	Language               *string    `json:"language,omitempty"`
	Status                 Status     `json:"status"`
	ProcessingStartedAt    *time.Time `json:"processing_started_at,omitempty"`
	ProcessingCompletedAt  *time.Time `json:"processing_completed_at,omitempty"`
	UploadedBy             *uuid.UUID `json:"uploaded_by,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
}

type Participant struct {
	ID               uuid.UUID `json:"id"`
	MeetingID        uuid.UUID `json:"meeting_id"`
	ParticipantName  string    `json:"participant_name"`
	ParticipantEmail *string   `json:"participant_email,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type TaskRelation struct {
	ID           uuid.UUID
	WorkspaceID  uuid.UUID
	ProjectID    *uuid.UUID
	MeetingID    *uuid.UUID
	Title        string
	Description  string
	Status       string
	Priority     string
	DueDate      *time.Time
	AIGenerated  bool
	AIConfidence float64
	AssignedTo   *uuid.UUID
}

type DecisionRelation struct {
	ID           uuid.UUID
	WorkspaceID  uuid.UUID
	ProjectID    *uuid.UUID
	MeetingID    uuid.UUID
	DecisionText string
}

type BlockerRelation struct {
	ID          uuid.UUID
	WorkspaceID uuid.UUID
	ProjectID   *uuid.UUID
	MeetingID   uuid.UUID
	BlockerText string
	Severity    string
	Resolved    bool
}

type DocumentRelation struct {
	ID            uuid.UUID
	WorkspaceID   uuid.UUID
	ProjectID     *uuid.UUID
	MeetingID     *uuid.UUID
	Title         string
	Type          string
	Content       string
	GeneratedByAI bool
}

type Repository interface {
	Create(ctx context.Context, m *Meeting) error
	GetByID(ctx context.Context, id uuid.UUID) (*Meeting, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*Meeting, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Meeting, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status Status) error
	Update(ctx context.Context, m *Meeting) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddParticipant(ctx context.Context, p *Participant) error
	ListParticipants(ctx context.Context, meetingID uuid.UUID) ([]*Participant, error)

	CreateTask(ctx context.Context, t *TaskRelation) error
	CreateDecision(ctx context.Context, d *DecisionRelation) error
	CreateBlocker(ctx context.Context, b *BlockerRelation) error
	CreateDocument(ctx context.Context, d *DocumentRelation) error
}

