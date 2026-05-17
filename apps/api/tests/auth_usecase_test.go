package auth_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/internal/domain/workspace"
	infraauth "github.com/meetext/backend/internal/infrastructure/auth"
	ucauth "github.com/meetext/backend/internal/usecase/auth"
	"github.com/meetext/backend/internal/config"
	"github.com/meetext/backend/pkg/apperr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"time"
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
func (m *mockUserRepo) Update(ctx context.Context, u *user.User) error {
	return m.Called(ctx, u).Error(0)
}
func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

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
func (m *mockWorkspaceRepo) GetBySlug(ctx context.Context, slug string) (*workspace.Workspace, error) {
	args := m.Called(ctx, slug)
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
func (m *mockWorkspaceRepo) UpdateMemberRole(ctx context.Context, wsID, uID uuid.UUID, role string) error {
	return m.Called(ctx, wsID, uID, role).Error(0)
}
func (m *mockWorkspaceRepo) RemoveMember(ctx context.Context, wsID, uID uuid.UUID) error {
	return m.Called(ctx, wsID, uID).Error(0)
}

// ── Tests ─────────────────────────────────────────────────────────────────────

func newTestUseCase(ur *mockUserRepo, wr *mockWorkspaceRepo) *ucauth.UseCase {
	jwtCfg := config.JWTConfig{
		AccessSecret:  "test-access-secret-32-chars-long!!",
		RefreshSecret: "test-refresh-secret-32-chars-long!",
		AccessTTL:     15 * time.Minute,
		RefreshTTL:    7 * 24 * time.Hour,
	}
	return ucauth.NewUseCase(ur, wr, infraauth.NewJWTService(jwtCfg))
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
	})

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
	})

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
	})

	assert.ErrorIs(t, err, apperr.ErrInvalidCredentials)
}
