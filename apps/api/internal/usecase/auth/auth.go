package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/internal/domain/workspace"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/internal/infrastructure/password"
	"github.com/meetext/backend/pkg/apperr"
)

type RegisterInput struct {
	FullName      string `json:"full_name"      validate:"required,min=2,max=100"`
	Email         string `json:"email"          validate:"required,email"`
	Password      string `json:"password"       validate:"required,min=8"`
	WorkspaceName string `json:"workspace_name" validate:"required,min=2,max=100"`
}

type LoginInput struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         *user.User           `json:"user"`
	Workspace    *workspace.Workspace `json:"workspace,omitempty"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type UseCase struct {
	userRepo      user.Repository
	workspaceRepo workspace.Repository
	jwt           *infraauth.JWTService
}

func NewUseCase(
	userRepo user.Repository,
	workspaceRepo workspace.Repository,
	jwt *infraauth.JWTService,
) *UseCase {
	return &UseCase{userRepo: userRepo, workspaceRepo: workspaceRepo, jwt: jwt}
}

func (uc *UseCase) Register(ctx context.Context, in RegisterInput) (*AuthResponse, error) {
	existing, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if existing != nil {
		return nil, apperr.ErrConflict
	}
	if err != nil && err != apperr.ErrNotFound {
		return nil, fmt.Errorf("auth: register: %w", err)
	}

	hash, err := password.Hash(in.Password)
	if err != nil {
		return nil, fmt.Errorf("auth: hash password: %w", err)
	}

	now := time.Now()
	u := &user.User{
		ID:           uuid.New(),
		FullName:     in.FullName,
		Email:        in.Email,
		PasswordHash: hash,
		Plan:         user.PlanFree,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("auth: create user: %w", err)
	}

	ws := &workspace.Workspace{
		ID:        uuid.New(),
		OwnerID:   u.ID,
		Name:      in.WorkspaceName,
		CreatedAt: now,
	}
	if err := uc.workspaceRepo.Create(ctx, ws); err != nil {
		return nil, fmt.Errorf("auth: create workspace: %w", err)
	}

	member := &workspace.Member{
		ID:          uuid.New(),
		WorkspaceID: ws.ID,
		UserID:      u.ID,
		Role:        workspace.RoleOwner,
		CreatedAt:   now,
	}
	if err := uc.workspaceRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("auth: add member: %w", err)
	}

	tokens, err := uc.jwt.IssueTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}

	return &AuthResponse{
		User:         u,
		Workspace:    ws,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (uc *UseCase) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	u, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, apperr.ErrInvalidCredentials
	}

	if !password.Compare(u.PasswordHash, in.Password) {
		return nil, apperr.ErrInvalidCredentials
	}

	tokens, err := uc.jwt.IssueTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}

	return &AuthResponse{
		User:         u,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	}, nil
}

func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (*infraauth.TokenPair, error) {
	claims, err := uc.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	if _, err := uc.userRepo.GetByID(ctx, claims.UserID); err != nil {
		return nil, apperr.ErrUnauthorized
	}
	return uc.jwt.IssueTokenPair(claims.UserID)
}
