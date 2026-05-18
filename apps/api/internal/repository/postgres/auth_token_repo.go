package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	authdomain "github.com/meetext/backend/internal/domain/auth"
	"github.com/meetext/backend/pkg/apperr"
)

type AuthTokenRepository struct {
	db *pgxpool.Pool
}

func NewAuthTokenRepository(db *pgxpool.Pool) *AuthTokenRepository {
	return &AuthTokenRepository{db: db}
}

func (r *AuthTokenRepository) CreateVerificationToken(ctx context.Context, token *authdomain.VerificationToken) error {
	q := `INSERT INTO verification_tokens (id, user_id, token_hash, expires_at, used_at, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6)`
	if _, err := r.db.Exec(ctx, q, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.UsedAt, token.CreatedAt); err != nil {
		return fmt.Errorf("auth token repo: create verification token: %w", err)
	}
	return nil
}

func (r *AuthTokenRepository) GetVerificationToken(ctx context.Context, tokenHash string) (*authdomain.VerificationToken, error) {
	q := `SELECT id, user_id, token_hash, expires_at, used_at, created_at
	      FROM verification_tokens WHERE token_hash=$1`
	token := &authdomain.VerificationToken{}
	err := r.db.QueryRow(ctx, q, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("auth token repo: get verification token: %w", err)
	}
	return token, nil
}

func (r *AuthTokenRepository) MarkVerificationTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	q := `UPDATE verification_tokens SET used_at=$1 WHERE id=$2 AND used_at IS NULL`
	tag, err := r.db.Exec(ctx, q, usedAt, id)
	if err != nil {
		return fmt.Errorf("auth token repo: mark verification token used: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *AuthTokenRepository) DeleteVerificationTokensForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM verification_tokens WHERE user_id=$1 AND used_at IS NULL`, userID)
	if err != nil {
		return fmt.Errorf("auth token repo: delete verification tokens: %w", err)
	}
	return nil
}

func (r *AuthTokenRepository) CreatePasswordResetToken(ctx context.Context, token *authdomain.PasswordResetToken) error {
	q := `INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at, used_at, created_at)
	      VALUES ($1,$2,$3,$4,$5,$6)`
	if _, err := r.db.Exec(ctx, q, token.ID, token.UserID, token.TokenHash, token.ExpiresAt, token.UsedAt, token.CreatedAt); err != nil {
		return fmt.Errorf("auth token repo: create password reset token: %w", err)
	}
	return nil
}

func (r *AuthTokenRepository) GetPasswordResetToken(ctx context.Context, tokenHash string) (*authdomain.PasswordResetToken, error) {
	q := `SELECT id, user_id, token_hash, expires_at, used_at, created_at
	      FROM password_reset_tokens WHERE token_hash=$1`
	token := &authdomain.PasswordResetToken{}
	err := r.db.QueryRow(ctx, q, tokenHash).Scan(
		&token.ID, &token.UserID, &token.TokenHash, &token.ExpiresAt, &token.UsedAt, &token.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("auth token repo: get password reset token: %w", err)
	}
	return token, nil
}

func (r *AuthTokenRepository) MarkPasswordResetTokenUsed(ctx context.Context, id uuid.UUID, usedAt time.Time) error {
	q := `UPDATE password_reset_tokens SET used_at=$1 WHERE id=$2 AND used_at IS NULL`
	tag, err := r.db.Exec(ctx, q, usedAt, id)
	if err != nil {
		return fmt.Errorf("auth token repo: mark password reset token used: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *AuthTokenRepository) DeletePasswordResetTokensForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM password_reset_tokens WHERE user_id=$1 AND used_at IS NULL`, userID)
	if err != nil {
		return fmt.Errorf("auth token repo: delete password reset tokens: %w", err)
	}
	return nil
}
