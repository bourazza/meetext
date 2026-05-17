package document

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeSummary      Type = "summary"
	TypeRequirements Type = "requirements"
	TypeTechnicalDoc Type = "technical_doc"
	TypeSprintPlan   Type = "sprint_plan"
	TypeClientNotes  Type = "client_notes"
	TypeDecisionLog  Type = "decision_log"
)

type Document struct {
	ID             uuid.UUID  `json:"id"`
	WorkspaceID    uuid.UUID  `json:"workspace_id"`
	ProjectID      uuid.UUID  `json:"project_id"`
	MeetingID      *uuid.UUID `json:"meeting_id,omitempty"`
	Title          string     `json:"title"`
	Type           Type       `json:"type"`
	Content        *string    `json:"content,omitempty"`
	GeneratedByAI  bool       `json:"generated_by_ai"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, d *Document) error
	GetByID(ctx context.Context, id uuid.UUID) (*Document, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Document, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Document, error)
	Update(ctx context.Context, d *Document) error
	Delete(ctx context.Context, id uuid.UUID) error
}
