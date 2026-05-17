package project

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPlanning   Status = "planning"
	StatusActive     Status = "active"
	StatusReview     Status = "review"
	StatusCompleted  Status = "completed"
)

type Project struct {
	ID          uuid.UUID  `json:"id"`
	WorkspaceID uuid.UUID  `json:"workspace_id"`
	ClientID    *uuid.UUID `json:"client_id,omitempty"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      Status     `json:"status"`
	Progress    int        `json:"progress"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Repository interface {
	Create(ctx context.Context, p *Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*Project, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*Project, error)
	Update(ctx context.Context, p *Project) error
	Delete(ctx context.Context, id uuid.UUID) error
}
