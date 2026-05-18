package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/internal/domain/workspace"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	infraoauth "github.com/meetext/backend/internal/infrastructure/oauth"
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
		Provider:     user.ProviderLocal,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("auth: create user: %w", err)
	}

	ws, err := uc.createWorkspace(ctx, u.ID, in.WorkspaceName, now)
	if err != nil {
		return nil, err
	}

	tokens, err := uc.jwt.IssueTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	return &AuthResponse{User: u, Workspace: ws, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (uc *UseCase) Login(ctx context.Context, in LoginInput) (*AuthResponse, error) {
	u, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, apperr.ErrInvalidCredentials
	}
	if u.Provider != user.ProviderLocal || u.PasswordHash == "" {
		return nil, apperr.ErrInvalidCredentials
	}
	if !password.Compare(u.PasswordHash, in.Password) {
		return nil, apperr.ErrInvalidCredentials
	}

	tokens, err := uc.jwt.IssueTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	return &AuthResponse{User: u, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
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

// OAuthLogin finds or creates a user from an OAuth provider profile, then issues a JWT pair.
func (uc *UseCase) OAuthLogin(ctx context.Context, provider user.Provider, info *infraoauth.UserInfo) (*AuthResponse, error) {
	// 1. Try to find by provider + provider_id (fastest path)
	u, err := uc.userRepo.GetByProviderID(ctx, provider, info.ProviderID)
	if err != nil && err != apperr.ErrNotFound {
		return nil, fmt.Errorf("auth: oauth lookup: %w", err)
	}

	if u == nil {
		// 2. Try to find by email — link existing account
		u, err = uc.userRepo.GetByEmail(ctx, info.Email)
		if err != nil && err != apperr.ErrNotFound {
			return nil, fmt.Errorf("auth: oauth email lookup: %w", err)
		}

		if u == nil {
			// 3. Create new user
			u, err = uc.createOAuthUser(ctx, provider, info)
			if err != nil {
				return nil, err
			}
		} else {
			// 4. Update existing user with provider info so future logins use path 1
			pid := info.ProviderID
			u.Provider = provider
			u.ProviderID = &pid
			u.UpdatedAt = time.Now()
			if u.AvatarURL == nil && info.AvatarURL != "" {
				u.AvatarURL = &info.AvatarURL
			}
			if err := uc.userRepo.Update(ctx, u); err != nil {
				return nil, fmt.Errorf("auth: oauth update user: %w", err)
			}
		}
	}

	// Fetch first workspace for response
	workspaces, _ := uc.workspaceRepo.ListByUserID(ctx, u.ID)
	var ws *workspace.Workspace
	if len(workspaces) > 0 {
		ws = workspaces[0]
	}

	tokens, err := uc.jwt.IssueTokenPair(u.ID)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	return &AuthResponse{User: u, Workspace: ws, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (uc *UseCase) createOAuthUser(ctx context.Context, provider user.Provider, info *infraoauth.UserInfo) (*user.User, error) {
	now := time.Now()
	pid := info.ProviderID

	var avatarURL *string
	if info.AvatarURL != "" {
		avatarURL = &info.AvatarURL
	}

	u := &user.User{
		ID:         uuid.New(),
		FullName:   info.Name,
		Email:      info.Email,
		AvatarURL:  avatarURL,
		Plan:       user.PlanFree,
		Provider:   provider,
		ProviderID: &pid,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("auth: create oauth user: %w", err)
	}

	// Auto-create workspace from name
	wsName := derivedWorkspaceName(info.Name)
	if _, err := uc.createWorkspace(ctx, u.ID, wsName, now); err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UseCase) createWorkspace(ctx context.Context, ownerID uuid.UUID, name string, now time.Time) (*workspace.Workspace, error) {
	ws := &workspace.Workspace{
		ID:        uuid.New(),
		OwnerID:   ownerID,
		Name:      name,
		CreatedAt: now,
	}
	if err := uc.workspaceRepo.Create(ctx, ws); err != nil {
		return nil, fmt.Errorf("auth: create workspace: %w", err)
	}
	member := &workspace.Member{
		ID:          uuid.New(),
		WorkspaceID: ws.ID,
		UserID:      ownerID,
		Role:        workspace.RoleOwner,
		CreatedAt:   now,
	}
	if err := uc.workspaceRepo.AddMember(ctx, member); err != nil {
		return nil, fmt.Errorf("auth: add member: %w", err)
	}
	return ws, nil
}

func derivedWorkspaceName(fullName string) string {
	name := strings.TrimSpace(fullName)
	if name == "" {
		return "My Workspace"
	}
	parts := strings.Fields(name)
	return parts[0] + "'s Workspace"
}
