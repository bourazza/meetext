package task

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusReview     Status = "review"
	StatusDone       Status = "done"
)

type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

type Task struct {
	ID           uuid.UUID  `json:"id"`
	WorkspaceID  uuid.UUID  `json:"workspace_id"`
	ProjectID    uuid.UUID  `json:"project_id"`
	MeetingID    *uuid.UUID `json:"meeting_id,omitempty"`
	Title        string     `json:"title"`
	Description  *string    `json:"description,omitempty"`
	Status       Status     `json:"status"`
	Priority     Priority   `json:"priority"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	AIGenerated  bool       `json:"ai_generated"`
	AIConfidence *float64   `json:"ai_confidence,omitempty"`
	AssignedTo   *uuid.UUID `json:"assigned_to,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, t *Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*Task, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]*Task, error)
	ListByMeeting(ctx context.Context, meetingID uuid.UUID) ([]*Task, error)
	Update(ctx context.Context, t *Task) error
	Delete(ctx context.Context, id uuid.UUID) error
}
