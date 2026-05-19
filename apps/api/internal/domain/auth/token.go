package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type VerificationToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PasswordResetToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type Session struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	RefreshHash string
	UserAgent   *string
	IPAddress   *string
	ExpiresAt   time.Time
	RevokedAt   *time.Time
	CreatedAt   time.Time
}

type TokenRepository interface {
	CreateVerificationToken(ctx context.Context, token *VerificationToken) error
	GetVerificationToken(ctx context.Context, tokenHash string) (*VerificationToken, error)
	MarkVerificationTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error
	DeleteVerificationTokensForUser(ctx context.Context, userID uuid.UUID) error

	CreatePasswordResetToken(ctx context.Context, token *PasswordResetToken) error
	GetPasswordResetToken(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error
	DeletePasswordResetTokensForUser(ctx context.Context, userID uuid.UUID) error

	CreateSession(ctx context.Context, session *Session) error
	GetSession(ctx context.Context, id uuid.UUID) (*Session, error)
	UpdateSessionRefreshHash(ctx context.Context, id uuid.UUID, refreshHash string, expiresAt time.Time) error
	RevokeSession(ctx context.Context, id uuid.UUID, revokedAt time.Time) error
	RevokeSessionsForUser(ctx context.Context, userID uuid.UUID, revokedAt time.Time) error
}
