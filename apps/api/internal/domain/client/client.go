package client

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID           uuid.UUID `json:"id"`
	WorkspaceID  uuid.UUID `json:"workspace_id"`
	CompanyName  string    `json:"company_name"`
	ContactName  *string   `json:"contact_name,omitempty"`
	ContactEmail *string   `json:"contact_email,omitempty"`
	LogoURL      *string   `json:"logo_url,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, c *Client) error
	GetByID(ctx context.Context, id uuid.UUID) (*Client, error)
	ListByWorkspace(ctx context.Context, workspaceID uuid.UUID) ([]*Client, error)
	Update(ctx context.Context, c *Client) error
	Delete(ctx context.Context, id uuid.UUID) error
}
