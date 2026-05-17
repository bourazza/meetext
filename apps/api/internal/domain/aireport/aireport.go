package aireport

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	Assignee    *string `json:"assignee,omitempty"`
	Priority    string  `json:"priority"` // low | medium | high
	DueDate     *string `json:"due_date,omitempty"`
}

type Decision struct {
	Description string  `json:"description"`
	MadeBy      *string `json:"made_by,omitempty"`
}

type Risk struct {
	Description string `json:"description"`
	Severity    string `json:"severity"` // low | medium | high
}

type AIReport struct {
	ID          uuid.UUID  `json:"id"`
	MeetingID   uuid.UUID  `json:"meeting_id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	Summary     string     `json:"summary"`
	Tasks       []Task     `json:"tasks"`
	Goals       []string   `json:"goals"`
	Decisions   []Decision `json:"decisions"`
	Risks       []Risk     `json:"risks"`
	Deadlines   []string   `json:"deadlines"`
	ModelUsed   string     `json:"model_used"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, r *AIReport) error
	GetByMeetingID(ctx context.Context, meetingID uuid.UUID) (*AIReport, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID, limit, offset int) ([]*AIReport, error)
	Update(ctx context.Context, r *AIReport) error
	Delete(ctx context.Context, meetingID uuid.UUID) error
}
