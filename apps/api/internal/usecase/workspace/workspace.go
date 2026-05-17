package workspace

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/workspace"
	"github.com/meetext/backend/pkg/apperr"
)

type UseCase struct {
	repo workspace.Repository
}

func NewUseCase(repo workspace.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) GetByID(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *UseCase) ListForUser(ctx context.Context, userID uuid.UUID) ([]*workspace.Workspace, error) {
	return uc.repo.ListByUserID(ctx, userID)
}

func (uc *UseCase) UpdateName(ctx context.Context, id uuid.UUID, name string, requesterID uuid.UUID) (*workspace.Workspace, error) {
	ws, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	member, err := uc.repo.GetMember(ctx, id, requesterID)
	if err != nil || (member.Role != workspace.RoleOwner && member.Role != workspace.RoleAdmin) {
		return nil, apperr.ErrForbidden
	}

	ws.Name = name
	if err := uc.repo.Update(ctx, ws); err != nil {
		return nil, fmt.Errorf("workspace: update: %w", err)
	}
	return ws, nil
}

func (uc *UseCase) ListMembers(ctx context.Context, workspaceID uuid.UUID) ([]*workspace.Member, error) {
	return uc.repo.ListMembers(ctx, workspaceID)
}

func (uc *UseCase) RemoveMember(ctx context.Context, workspaceID, userID, requesterID uuid.UUID) error {
	requester, err := uc.repo.GetMember(ctx, workspaceID, requesterID)
	if err != nil || (requester.Role != workspace.RoleOwner && requester.Role != workspace.RoleAdmin) {
		return apperr.ErrForbidden
	}
	return uc.repo.RemoveMember(ctx, workspaceID, userID)
}
