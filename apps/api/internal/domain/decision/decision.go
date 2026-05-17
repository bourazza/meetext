package decision

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Decision struct {
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	ProjectID    uuid.UUID `json:"project_id"`
	MeetingID    uuid.UUID `json:"meeting_id"`
	DecisionText string    `json:"decision_text"`
	CreatedAt    time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, d *Decision) error
	GetByID(ctx context.Context, id uuid.UUID) (*Decision, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Decision, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Decision, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
