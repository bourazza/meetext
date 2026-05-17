package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meetext/backend/internal/domain/workspace"
	"github.com/meetext/backend/pkg/apperr"
)

type WorkspaceRepository struct {
	db *pgxpool.Pool
}

func NewWorkspaceRepository(db *pgxpool.Pool) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, w *workspace.Workspace) error {
	q := `INSERT INTO workspaces (id, owner_id, name, logo_url, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, w.ID, w.OwnerID, w.Name, w.LogoURL, w.CreatedAt)
	if err != nil {
		return fmt.Errorf("workspace repo: create: %w", err)
	}
	return nil
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	q := `SELECT id, owner_id, name, logo_url, created_at FROM workspaces WHERE id=$1`
	w := &workspace.Workspace{}
	err := r.db.QueryRow(ctx, q, id).Scan(&w.ID, &w.OwnerID, &w.Name, &w.LogoURL, &w.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("workspace repo: get by id: %w", err)
	}
	return w, nil
}

func (r *WorkspaceRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*workspace.Workspace, error) {
	q := `SELECT w.id, w.owner_id, w.name, w.logo_url, w.created_at
		  FROM workspaces w
		  JOIN workspace_members wm ON wm.workspace_id = w.id
		  WHERE wm.user_id = $1`
	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("workspace repo: list by user: %w", err)
	}
	defer rows.Close()

	var result []*workspace.Workspace
	for rows.Next() {
		w := &workspace.Workspace{}
		if err := rows.Scan(&w.ID, &w.OwnerID, &w.Name, &w.LogoURL, &w.CreatedAt); err != nil {
			return nil, fmt.Errorf("workspace repo: scan: %w", err)
		}
		result = append(result, w)
	}
	return result, nil
}

func (r *WorkspaceRepository) Update(ctx context.Context, w *workspace.Workspace) error {
	q := `UPDATE workspaces SET name=$1, logo_url=$2 WHERE id=$3`
	tag, err := r.db.Exec(ctx, q, w.Name, w.LogoURL, w.ID)
	if err != nil {
		return fmt.Errorf("workspace repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *WorkspaceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM workspaces WHERE id=$1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("workspace repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *WorkspaceRepository) AddMember(ctx context.Context, m *workspace.Member) error {
	q := `INSERT INTO workspace_members (id, workspace_id, user_id, role, created_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := r.db.Exec(ctx, q, m.ID, m.WorkspaceID, m.UserID, m.Role, m.CreatedAt)
	if err != nil {
		if isDuplicateKey(err) {
			return apperr.ErrConflict
		}
		return fmt.Errorf("workspace repo: add member: %w", err)
	}
	return nil
}

func (r *WorkspaceRepository) GetMember(ctx context.Context, workspaceID, userID uuid.UUID) (*workspace.Member, error) {
	q := `SELECT id, workspace_id, user_id, role, created_at FROM workspace_members WHERE workspace_id=$1 AND user_id=$2`
	m := &workspace.Member{}
	err := r.db.QueryRow(ctx, q, workspaceID, userID).Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("workspace repo: get member: %w", err)
	}
	return m, nil
}

func (r *WorkspaceRepository) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	q := `SELECT id, workspace_id, user_id, role, created_at FROM workspace_members WHERE workspace_id=$1`
	rows, err := r.db.Query(ctx, q, workspaceID)
	if err != nil {
		return nil, fmt.Errorf("workspace repo: list members: %w", err)
	}
	defer rows.Close()

	var result []*workspace.Member
	for rows.Next() {
		m := &workspace.Member{}
		if err := rows.Scan(&m.ID, &m.WorkspaceID, &m.UserID, &m.Role, &m.CreatedAt); err != nil {
			return nil, fmt.Errorf("workspace repo: scan member: %w", err)
		}
		result = append(result, m)
	}
	return result, nil
}

func (r *WorkspaceRepository) UpdateMemberRole(ctx context.Context, workspaceID, userID uuid.UUID, role workspace.Role) error {
	q := `UPDATE workspace_members SET role=$1 WHERE workspace_id=$2 AND user_id=$3`
	tag, err := r.db.Exec(ctx, q, role, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("workspace repo: update member role: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *WorkspaceRepository) RemoveMember(ctx context.Context, workspaceID, userID uuid.UUID) error {
	q := `DELETE FROM workspace_members WHERE workspace_id=$1 AND user_id=$2`
	tag, err := r.db.Exec(ctx, q, workspaceID, userID)
	if err != nil {
		return fmt.Errorf("workspace repo: remove member: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
