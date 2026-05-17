package goal

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Goal struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	MeetingID   *uuid.UUID `json:"meeting_id,omitempty"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Completed   bool       `json:"completed"`
	TargetDate  *time.Time `json:"target_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, g *Goal) error
	GetByID(ctx context.Context, id uuid.UUID) (*Goal, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Goal, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Goal, error)
	Update(ctx context.Context, g *Goal) error
	Delete(ctx context.Context, id uuid.UUID) error
}
