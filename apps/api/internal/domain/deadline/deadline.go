package deadline

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Deadline struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ProjectID   uuid.UUID  `json:"project_id"`
	MeetingID   *uuid.UUID `json:"meeting_id,omitempty"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	DueDate     time.Time  `json:"due_date"`
	Completed   bool       `json:"completed"`
	CreatedAt   time.Time  `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, d *Deadline) error
	GetByID(ctx context.Context, id uuid.UUID) (*Deadline, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Deadline, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Deadline, error)
	Update(ctx context.Context, d *Deadline) error
	Delete(ctx context.Context, id uuid.UUID) error
}
