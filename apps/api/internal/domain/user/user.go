package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Plan string

const (
	PlanFree     Plan = "free"
	PlanPro      Plan = "pro"
	PlanBusiness Plan = "business"
)

type Provider string

const (
	ProviderLocal  Provider = "local"
	ProviderGoogle Provider = "google"
	ProviderGitHub Provider = "github"
)

type User struct {
	ID              uuid.UUID  `json:"id"`
	FullName        string     `json:"full_name"`
	Email           string     `json:"email"`
	PasswordHash    string     `json:"-"`
	AvatarURL       *string    `json:"avatar_url,omitempty"`
	Plan            Plan       `json:"plan"`
	Provider        Provider   `json:"provider"`
	ProviderID      *string    `json:"-"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type OAuthAccount struct {
	ID                uuid.UUID
	UserID            uuid.UUID
	Provider          Provider
	ProviderAccountID string
	Email             string
	AvatarURL         *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type Repository interface {
	Create(ctx context.Context, u *User) error
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByProviderID(ctx context.Context, provider Provider, providerID string) (*User, error)
	GetByOAuthAccount(ctx context.Context, provider Provider, providerAccountID string) (*User, error)
	UpsertOAuthAccount(ctx context.Context, account *OAuthAccount) error
	Update(ctx context.Context, u *User) error
	MarkEmailVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string, updatedAt time.Time) error
	RecordLogin(ctx context.Context, id uuid.UUID, loggedInAt time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}
