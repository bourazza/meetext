package workspace

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type Workspace struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	Name      string    `json:"name"`
	LogoURL   *string   `json:"logo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Member struct {
	ID          uuid.UUID `json:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id"`
	UserID      uuid.UUID `json:"user_id"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repository interface {
	Create(ctx context.Context, w *Workspace) error
	GetByID(ctx context.Context, id uuid.UUID) (*Workspace, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]*Workspace, error)
	Update(ctx context.Context, w *Workspace) error
	Delete(ctx context.Context, id uuid.UUID) error

	AddMember(ctx context.Context, m *Member) error
	GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*Member, error)
	ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*Member, error)
	UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role Role) error
	RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error
}
