package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	authdomain "github.com/meetext/backend/internal/domain/auth"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/internal/domain/workspace"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	"github.com/meetext/backend/internal/infrastructure/email"
	infraoauth "github.com/meetext/backend/internal/infrastructure/oauth"
	"github.com/meetext/backend/internal/infrastructure/password"
	"github.com/meetext/backend/pkg/apperr"
)

type RegisterInput struct {
	FullName      string `json:"full_name"      validate:"required,min=2,max=100"`
	Email         string `json:"email"          validate:"required,email"`
	Password      string `json:"password"       validate:"required,min=8"`
	WorkspaceName string `json:"workspace_name" validate:"omitempty,min=2,max=100"`
}

type LoginInput struct {
	Email      string `json:"email"       validate:"required,email"`
	Password   string `json:"password"    validate:"required"`
	RememberMe bool   `json:"remember_me"`
}

type ForgotPasswordInput struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordInput struct {
	Token    string `json:"token"    validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}

type VerifyEmailInput struct {
	Token string `json:"token" validate:"required"`
}

type ResendVerificationInput struct {
	Email string `json:"email" validate:"required,email"`
}

type RequestMeta struct {
	UserAgent string
	IP        string
	Remember  bool
}

type AuthResponse struct {
	User         *user.User           `json:"user"`
	Workspace    *workspace.Workspace `json:"workspace,omitempty"`
	AccessToken  string               `json:"access_token"`
	RefreshToken string               `json:"refresh_token"`
}

type UseCase struct {
	userRepo        user.Repository
	workspaceRepo   workspace.Repository
	tokenRepo       authdomain.TokenRepository
	jwt             *infraauth.JWTService
	email           email.Service
	frontendURL     string
	requireVerified bool
}

func NewUseCase(
	userRepo user.Repository,
	workspaceRepo workspace.Repository,
	tokenRepo authdomain.TokenRepository,
	jwt *infraauth.JWTService,
	email email.Service,
	frontendURL string,
	requireVerified bool,
) *UseCase {
	return &UseCase{
		userRepo:        userRepo,
		workspaceRepo:   workspaceRepo,
		tokenRepo:       tokenRepo,
		jwt:             jwt,
		email:           email,
		frontendURL:     strings.TrimRight(frontendURL, "/"),
		requireVerified: requireVerified,
	}
}

func (uc *UseCase) Register(ctx context.Context, in RegisterInput, meta RequestMeta) (*AuthResponse, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))
	in.FullName = strings.TrimSpace(in.FullName)

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

	workspaceName := strings.TrimSpace(in.WorkspaceName)
	if workspaceName == "" {
		workspaceName = derivedWorkspaceName(in.FullName)
	}
	ws, err := uc.createWorkspace(ctx, u.ID, workspaceName, now)
	if err != nil {
		return nil, err
	}

	if err := uc.issueVerificationEmail(ctx, u); err != nil {
		return nil, err
	}

	tokens, err := uc.createSessionAndTokens(ctx, u.ID, meta)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	return &AuthResponse{User: u, Workspace: ws, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (uc *UseCase) Login(ctx context.Context, in LoginInput, meta RequestMeta) (*AuthResponse, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	u, err := uc.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, apperr.ErrInvalidCredentials
	}
	if u.PasswordHash == "" {
		return nil, apperr.ErrInvalidCredentials
	}
	if !password.Compare(u.PasswordHash, in.Password) {
		return nil, apperr.ErrInvalidCredentials
	}
	if uc.requireVerified && u.EmailVerifiedAt == nil {
		return nil, apperr.ErrEmailNotVerified
	}

	tokens, err := uc.createSessionAndTokens(ctx, u.ID, meta)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	now := time.Now()
	_ = uc.userRepo.RecordLogin(ctx, u.ID, now)
	u.LastLoginAt = &now

	workspaces, _ := uc.workspaceRepo.ListByUserID(ctx, u.ID)
	var ws *workspace.Workspace
	if len(workspaces) > 0 {
		ws = workspaces[0]
	}
	return &AuthResponse{User: u, Workspace: ws, AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken}, nil
}

func (uc *UseCase) RefreshToken(ctx context.Context, refreshToken string) (*infraauth.TokenPair, error) {
	claims, err := uc.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	session, err := uc.tokenRepo.GetSession(ctx, claims.SessionID)
	if err != nil {
		return nil, apperr.ErrUnauthorized
	}
	now := time.Now()
	if session.UserID != claims.UserID || session.RevokedAt != nil || now.After(session.ExpiresAt) {
		return nil, apperr.ErrUnauthorized
	}
	if session.RefreshHash != hashToken(refreshToken) {
		_ = uc.tokenRepo.RevokeSession(ctx, session.ID, now)
		return nil, apperr.ErrUnauthorized
	}
	if _, err := uc.userRepo.GetByID(ctx, claims.UserID); err != nil {
		return nil, apperr.ErrUnauthorized
	}
	tokens, err := uc.jwt.IssueTokenPair(claims.UserID, claims.SessionID)
	if err != nil {
		return nil, err
	}
	if err := uc.tokenRepo.UpdateSessionRefreshHash(ctx, claims.SessionID, hashToken(tokens.RefreshToken), now.Add(uc.jwt.RefreshTTL())); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (uc *UseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	claims, err := uc.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil
	}
	_ = uc.tokenRepo.RevokeSession(ctx, claims.SessionID, time.Now())
	return nil
}

func (uc *UseCase) CurrentUser(ctx context.Context, userID uuid.UUID) (*user.User, error) {
	u, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperr.ErrUnauthorized
	}
	return u, nil
}

func (uc *UseCase) RequestPasswordReset(ctx context.Context, in ForgotPasswordInput) error {
	emailAddress := strings.ToLower(strings.TrimSpace(in.Email))
	u, err := uc.userRepo.GetByEmail(ctx, emailAddress)
	if err != nil {
		if err == apperr.ErrNotFound {
			return nil
		}
		return fmt.Errorf("auth: password reset lookup: %w", err)
	}

	raw, hashed, err := newToken()
	if err != nil {
		return fmt.Errorf("auth: password reset token: %w", err)
	}
	now := time.Now()
	if err := uc.tokenRepo.DeletePasswordResetTokensForUser(ctx, u.ID); err != nil {
		return err
	}
	if err := uc.tokenRepo.CreatePasswordResetToken(ctx, &authdomain.PasswordResetToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		TokenHash: hashed,
		ExpiresAt: now.Add(30 * time.Minute),
		CreatedAt: now,
	}); err != nil {
		return err
	}

	return uc.email.SendPasswordReset(ctx, u.Email, u.FullName, uc.frontendLink("/reset-password", raw))
}

func (uc *UseCase) ResetPassword(ctx context.Context, in ResetPasswordInput) error {
	hashed := hashToken(in.Token)
	token, err := uc.tokenRepo.GetPasswordResetToken(ctx, hashed)
	if err != nil {
		return apperr.New(400, "INVALID_RESET_TOKEN", "This password reset link is invalid or expired")
	}
	now := time.Now()
	if token.UsedAt != nil || now.After(token.ExpiresAt) {
		return apperr.New(400, "INVALID_RESET_TOKEN", "This password reset link is invalid or expired")
	}

	passwordHash, err := password.Hash(in.Password)
	if err != nil {
		return fmt.Errorf("auth: hash reset password: %w", err)
	}
	if err := uc.userRepo.UpdatePassword(ctx, token.UserID, passwordHash, now); err != nil {
		return err
	}
	if err := uc.tokenRepo.MarkPasswordResetTokenUsed(ctx, token.ID, now); err != nil {
		return err
	}
	_ = uc.tokenRepo.RevokeSessionsForUser(ctx, token.UserID, now)
	return nil
}

func (uc *UseCase) VerifyEmail(ctx context.Context, in VerifyEmailInput) error {
	hashed := hashToken(in.Token)
	token, err := uc.tokenRepo.GetVerificationToken(ctx, hashed)
	if err != nil {
		return apperr.New(400, "INVALID_VERIFICATION_TOKEN", "This verification link is invalid or expired")
	}
	now := time.Now()
	if token.UsedAt != nil || now.After(token.ExpiresAt) {
		return apperr.New(400, "INVALID_VERIFICATION_TOKEN", "This verification link is invalid or expired")
	}
	if err := uc.userRepo.MarkEmailVerified(ctx, token.UserID, now); err != nil {
		return err
	}
	if err := uc.tokenRepo.MarkVerificationTokenUsed(ctx, token.ID, now); err != nil {
		return err
	}
	return nil
}

func (uc *UseCase) ResendVerification(ctx context.Context, in ResendVerificationInput) error {
	emailAddress := strings.ToLower(strings.TrimSpace(in.Email))
	u, err := uc.userRepo.GetByEmail(ctx, emailAddress)
	if err != nil {
		if err == apperr.ErrNotFound {
			return nil
		}
		return fmt.Errorf("auth: resend verification lookup: %w", err)
	}
	if u.EmailVerifiedAt != nil {
		return nil
	}
	return uc.issueVerificationEmail(ctx, u)
}

// OAuthLogin finds or creates a user from an OAuth provider profile, then issues a JWT pair.
func (uc *UseCase) OAuthLogin(ctx context.Context, provider user.Provider, info *infraoauth.UserInfo, meta RequestMeta) (*AuthResponse, error) {
	// 1. Try to find by provider + provider_id (fastest path)
	u, err := uc.userRepo.GetByOAuthAccount(ctx, provider, info.ProviderID)
	if err != nil && err != apperr.ErrNotFound {
		return nil, fmt.Errorf("auth: oauth lookup: %w", err)
	}
	if u == nil {
		u, err = uc.userRepo.GetByProviderID(ctx, provider, info.ProviderID)
		if err != nil && err != apperr.ErrNotFound {
			return nil, fmt.Errorf("auth: oauth legacy lookup: %w", err)
		}
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
			now := time.Now()
			u.EmailVerifiedAt = &now
			if err := uc.userRepo.Update(ctx, u); err != nil {
				return nil, fmt.Errorf("auth: oauth update user: %w", err)
			}
		}
	}
	if err := uc.upsertOAuthAccount(ctx, u.ID, provider, info); err != nil {
		return nil, err
	}

	// Fetch first workspace for response
	workspaces, _ := uc.workspaceRepo.ListByUserID(ctx, u.ID)
	var ws *workspace.Workspace
	if len(workspaces) > 0 {
		ws = workspaces[0]
	}

	tokens, err := uc.createSessionAndTokens(ctx, u.ID, meta)
	if err != nil {
		return nil, fmt.Errorf("auth: issue tokens: %w", err)
	}
	now := time.Now()
	_ = uc.userRepo.RecordLogin(ctx, u.ID, now)
	u.LastLoginAt = &now
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
		ID:              uuid.New(),
		FullName:        info.Name,
		Email:           strings.ToLower(strings.TrimSpace(info.Email)),
		AvatarURL:       avatarURL,
		Plan:            user.PlanFree,
		Provider:        provider,
		ProviderID:      &pid,
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := uc.userRepo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("auth: create oauth user: %w", err)
	}
	if err := uc.upsertOAuthAccount(ctx, u.ID, provider, info); err != nil {
		return nil, err
	}

	// Auto-create workspace from name
	wsName := derivedWorkspaceName(info.Name)
	if _, err := uc.createWorkspace(ctx, u.ID, wsName, now); err != nil {
		return nil, err
	}
	return u, nil
}

func (uc *UseCase) upsertOAuthAccount(ctx context.Context, userID uuid.UUID, provider user.Provider, info *infraoauth.UserInfo) error {
	now := time.Now()
	var avatarURL *string
	if info.AvatarURL != "" {
		avatarURL = &info.AvatarURL
	}
	if err := uc.userRepo.UpsertOAuthAccount(ctx, &user.OAuthAccount{
		ID:                uuid.New(),
		UserID:            userID,
		Provider:          provider,
		ProviderAccountID: info.ProviderID,
		Email:             strings.ToLower(strings.TrimSpace(info.Email)),
		AvatarURL:         avatarURL,
		CreatedAt:         now,
		UpdatedAt:         now,
	}); err != nil {
		return fmt.Errorf("auth: upsert oauth account: %w", err)
	}
	return nil
}

func (uc *UseCase) createSessionAndTokens(ctx context.Context, userID uuid.UUID, meta RequestMeta) (*infraauth.TokenPair, error) {
	now := time.Now()
	sessionID := uuid.New()
	tokens, err := uc.jwt.IssueTokenPair(userID, sessionID)
	if err != nil {
		return nil, err
	}

	userAgent := nullableMeta(meta.UserAgent)
	ip := nullableMeta(meta.IP)
	refreshTTL := uc.jwt.RefreshTTL()
	if !meta.Remember {
		refreshTTL = 24 * time.Hour
	}
	session := &authdomain.Session{
		ID:          sessionID,
		UserID:      userID,
		RefreshHash: hashToken(tokens.RefreshToken),
		UserAgent:   userAgent,
		IPAddress:   ip,
		ExpiresAt:   now.Add(refreshTTL),
		CreatedAt:   now,
	}
	if err := uc.tokenRepo.CreateSession(ctx, session); err != nil {
		return nil, err
	}
	return tokens, nil
}

func nullableMeta(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func (uc *UseCase) issueVerificationEmail(ctx context.Context, u *user.User) error {
	raw, hashed, err := newToken()
	if err != nil {
		return fmt.Errorf("auth: verification token: %w", err)
	}
	now := time.Now()
	if err := uc.tokenRepo.DeleteVerificationTokensForUser(ctx, u.ID); err != nil {
		return err
	}
	if err := uc.tokenRepo.CreateVerificationToken(ctx, &authdomain.VerificationToken{
		ID:        uuid.New(),
		UserID:    u.ID,
		TokenHash: hashed,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}); err != nil {
		return err
	}
	return uc.email.SendVerification(ctx, u.Email, u.FullName, uc.frontendLink("/verify-email", raw))
}

func (uc *UseCase) frontendLink(path string, token string) string {
	u, err := url.Parse(uc.frontendURL + path)
	if err != nil {
		return uc.frontendURL + path + "?token=" + url.QueryEscape(token)
	}
	q := u.Query()
	q.Set("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}

func newToken() (raw string, hashed string, err error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", err
	}
	raw = base64.RawURLEncoding.EncodeToString(b)
	return raw, hashToken(raw), nil
}

func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
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
