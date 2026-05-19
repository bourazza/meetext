package auth_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/config"
	authdomain "github.com/meetext/backend/internal/domain/auth"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/internal/domain/workspace"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ── Mocks ────────────────────────────────────────────────────────────────────

type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) Create(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockUserRepo) GetByProviderID(ctx context.Context, provider user.Provider, providerID string) (*user.User, error) {
	args := m.Called(ctx, provider, providerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockUserRepo) GetByOAuthAccount(ctx context.Context, provider user.Provider, providerAccountID string) (*user.User, error) {
	args := m.Called(ctx, provider, providerAccountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}
func (m *mockUserRepo) UpsertOAuthAccount(ctx context.Context, account *user.OAuthAccount) error {
	return m.Called(ctx, account).Error(0)
}
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) MarkEmailVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error {
	return m.Called(ctx, id, verifiedAt).Error(0)
}
func (m *mockUserRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string, updatedAt time.Time) error {
	return m.Called(ctx, id, passwordHash, updatedAt).Error(0)
}
func (m *mockUserRepo) RecordLogin(ctx context.Context, id uuid.UUID, loggedInAt time.Time) error {
	return m.Called(ctx, id, loggedInAt).Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockTokenRepo struct{ mock.Mock }

func (m *mockTokenRepo) CreateVerificationToken(ctx context.Context, token *authdomain.VerificationToken) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockTokenRepo) GetVerificationToken(ctx context.Context, tokenHash string) (*authdomain.VerificationToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.VerificationToken), args.Error(1)
}
func (m *mockTokenRepo) MarkVerificationTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	return m.Called(ctx, id, usedAt).Error(0)
}
func (m *mockTokenRepo) DeleteVerificationTokensForUser(ctx context.Context, userID uuid.UUID) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *mockTokenRepo) CreatePasswordResetToken(ctx context.Context, token *authdomain.PasswordResetToken) error {
	return m.Called(ctx, token).Error(0)
}
func (m *mockTokenRepo) GetPasswordResetToken(ctx context.Context, tokenHash string) (*authdomain.PasswordResetToken, error) {
	args := m.Called(ctx, tokenHash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.PasswordResetToken), args.Error(1)
}
func (m *mockTokenRepo) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	return m.Called(ctx, id, usedAt).Error(0)
}
func (m *mockTokenRepo) DeletePasswordResetTokensForUser(ctx context.Context, userID uuid.UUID) error {
	return m.Called(ctx, userID).Error(0)
}
func (m *mockTokenRepo) CreateSession(ctx context.Context, session *authdomain.Session) error {
	return m.Called(ctx, session).Error(0)
}
func (m *mockTokenRepo) GetSession(ctx context.Context, id uuid.UUID) (*authdomain.Session, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authdomain.Session), args.Error(1)
}
func (m *mockTokenRepo) UpdateSessionRefreshHash(ctx context.Context, id uuid.UUID, refreshHash string, expiresAt time.Time) error {
	return m.Called(ctx, id, refreshHash, expiresAt).Error(0)
}
func (m *mockTokenRepo) RevokeSession(ctx context.Context, id uuid.UUID, revokedAt time.Time) error {
	return m.Called(ctx, id, revokedAt).Error(0)
}
func (m *mockTokenRepo) RevokeSessionsForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error {
	return m.Called(ctx, userID, revokedAt).Error(0)
}

type noopEmail struct{}

func (noopEmail) SendVerification(ctx context.Context, to, name, link string) error  { return nil }
func (noopEmail) SendPasswordReset(ctx context.Context, to, name, link string) error { return nil }

type mockWorkspaceRepo struct{ mock.Mock }

func (m *mockWorkspaceRepo) Create(ctx context.Context, w *workspace.Workspace) error {
	return m.Called(ctx, w).Error(0)
}
func (m *mockWorkspaceRepo) GetByID(ctx context.Context, id uuid.UUID) (*workspace.Workspace, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workspace.Workspace), args.Error(1)
}
func (m *mockWorkspaceRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*workspace.Workspace, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*workspace.Workspace), args.Error(1)
}
func (m *mockWorkspaceRepo) Update(ctx context.Context, w *workspace.Workspace) error {
	return m.Called(ctx, w).Error(0)
}
func (m *mockWorkspaceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockWorkspaceRepo) AddMember(ctx context.Context, mem *workspace.Member) error {
	return m.Called(ctx, mem).Error(0)
}
func (m *mockWorkspaceRepo) GetMember(ctx context.Context, wsID, uID uuid.UUID) (*workspace.Member, error) {
	args := m.Called(ctx, wsID, uID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workspace.Member), args.Error(1)
}
func (m *mockWorkspaceRepo) ListMembers(ctx context.Context, wsID uuid.UUID) ([]*workspace.Member, error) {
	args := m.Called(ctx, wsID)
	return args.Get(0).([]*workspace.Member), args.Error(1)
}
func (m *mockWorkspaceRepo) UpdateMemberRole(ctx context.Context, wsID, uID uuid.UUID, role workspace.Role) error {
	return m.Called(ctx, wsID, uID, role).Error(0)
}
func (m *mockWorkspaceRepo) RemoveMember(ctx context.Context, wsID, uID uuid.UUID) error {
	return m.Called(ctx, wsID, uID).Error(0)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func newTestUseCase(ur *mockUserRepo, wr *mockWorkspaceRepo) *ucauth.UseCase {
	tr := &mockTokenRepo{}
	jwtCfg := config.JWTConfig{
		AccessSecret:  "test-access-secret-32-chars-long!!",
		RefreshSecret: "test-refresh-secret-32-chars-long!",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}
	tr.On("DeleteVerificationTokensForUser", mock.Anything, mock.Anything).Maybe().Return(nil)
	tr.On("CreateVerificationToken", mock.Anything, mock.AnythingOfType("*auth.VerificationToken")).Maybe().Return(nil)
	tr.On("CreateSession", mock.Anything, mock.AnythingOfType("*auth.Session")).Maybe().Return(nil)
	return ucauth.NewUseCase(ur, wr, tr, infraauth.NewJWTService(jwtCfg), noopEmail{}, "http://localhost:3000", false)
}

func TestRegister_Success(t *testing.T) {
	ur := &mockUserRepo{}
	wr := &mockWorkspaceRepo{}
	uc := newTestUseCase(ur, wr)

	ur.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, apperr.ErrNotFound)
	ur.On("Create", mock.Anything, mock.AnythingOfType("*user.User")).Return(nil)
	wr.On("Create", mock.Anything, mock.AnythingOfType("*workspace.Workspace")).Return(nil)
	wr.On("AddMember", mock.Anything, mock.AnythingOfType("*workspace.Member")).Return(nil)

	res, err := uc.Register(context.Background(), ucauth.RegisterInput{
		FullName:      "Test User",
		Email:         "test@example.com",
		Password:      "securepassword",
		WorkspaceName: "Test Workspace",
	}, ucauth.RequestMeta{Remember: true})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.NotEmpty(t, res.AccessToken)
	assert.NotEmpty(t, res.RefreshToken)
	assert.Equal(t, "test@example.com", res.User.Email)
	ur.AssertExpectations(t)
	wr.AssertExpectations(t)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	ur := &mockUserRepo{}
	wr := &mockWorkspaceRepo{}
	uc := newTestUseCase(ur, wr)

	existing := &user.User{ID: uuid.New(), Email: "test@example.com"}
	ur.On("GetByEmail", mock.Anything, "test@example.com").Return(existing, nil)

	_, err := uc.Register(context.Background(), ucauth.RegisterInput{
		FullName:      "Test User",
		Email:         "test@example.com",
		Password:      "securepassword",
		WorkspaceName: "Test Workspace",
	}, ucauth.RequestMeta{Remember: true})

	assert.ErrorIs(t, err, apperr.ErrConflict)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	ur := &mockUserRepo{}
	wr := &mockWorkspaceRepo{}
	uc := newTestUseCase(ur, wr)

	ur.On("GetByEmail", mock.Anything, "wrong@example.com").Return(nil, apperr.ErrNotFound)

	_, err := uc.Login(context.Background(), ucauth.LoginInput{
		Email:    "wrong@example.com",
		Password: "anypassword",
	}, ucauth.RequestMeta{Remember: true})

	assert.ErrorIs(t, err, apperr.ErrInvalidCredentials)
}
