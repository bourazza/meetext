package blocker

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Blocker struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	MeetingID   uuid.UUID `json:"meeting_id"`
	BlockerText string    `json:"blocker_text"`
	Severity    *string   `json:"severity,omitempty"`
	Resolved    bool      `json:"resolved"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, b *Blocker) error
	GetByID(ctx context.Context, id uuid.UUID) (*Blocker, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Blocker, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Blocker, error)
	Update(ctx context.Context, b *Blocker) error
	Delete(ctx context.Context, id uuid.UUID) error
}
