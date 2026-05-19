package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/meetext/backend/internal/domain/user"
	"github.com/meetext/backend/pkg/apperr"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	q := `INSERT INTO users
		  (id, full_name, email, password_hash, avatar_url, plan, provider, provider_id, email_verified_at, last_login_at, created_at, updated_at)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	_, err := r.db.Exec(ctx, q,
		u.ID, u.FullName, u.Email, nullableString(u.PasswordHash),
		u.AvatarURL, u.Plan, u.Provider, u.ProviderID,
		u.EmailVerifiedAt, u.LastLoginAt,
		u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		if isDuplicateKey(err) {
			return apperr.ErrConflict
		}
		return fmt.Errorf("user repo: create: %w", err)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.User, error) {
	q := `SELECT id, full_name, email, COALESCE(password_hash,''), avatar_url, plan,
		         provider, provider_id, email_verified_at, last_login_at, created_at, updated_at
		  FROM users WHERE id=$1`
	return r.scanOne(r.db.QueryRow(ctx, q, id))
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	q := `SELECT id, full_name, email, COALESCE(password_hash,''), avatar_url, plan,
		         provider, provider_id, email_verified_at, last_login_at, created_at, updated_at
		  FROM users WHERE email=$1`
	return r.scanOne(r.db.QueryRow(ctx, q, email))
}

func (r *UserRepository) GetByProviderID(ctx context.Context, provider user.Provider, providerID string) (*user.User, error) {
	q := `SELECT id, full_name, email, COALESCE(password_hash,''), avatar_url, plan,
		         provider, provider_id, email_verified_at, last_login_at, created_at, updated_at
		  FROM users WHERE provider=$1 AND provider_id=$2`
	return r.scanOne(r.db.QueryRow(ctx, q, provider, providerID))
}

func (r *UserRepository) GetByOAuthAccount(ctx context.Context, provider user.Provider, providerAccountID string) (*user.User, error) {
	q := `SELECT u.id, u.full_name, u.email, COALESCE(u.password_hash,''), u.avatar_url, u.plan,
		         u.provider, u.provider_id, u.email_verified_at, u.last_login_at, u.created_at, u.updated_at
		  FROM users u
		  JOIN oauth_accounts oa ON oa.user_id = u.id
		  WHERE oa.provider=$1 AND oa.provider_account_id=$2`
	return r.scanOne(r.db.QueryRow(ctx, q, provider, providerAccountID))
}

func (r *UserRepository) UpsertOAuthAccount(ctx context.Context, account *user.OAuthAccount) error {
	q := `INSERT INTO oauth_accounts
		  (id, user_id, provider, provider_account_id, email, avatar_url, created_at, updated_at)
		  VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		  ON CONFLICT (provider, provider_account_id)
		  DO UPDATE SET user_id=EXCLUDED.user_id,
		                email=EXCLUDED.email,
		                avatar_url=EXCLUDED.avatar_url,
		                updated_at=EXCLUDED.updated_at`
	_, err := r.db.Exec(ctx, q,
		account.ID, account.UserID, account.Provider, account.ProviderAccountID,
		account.Email, account.AvatarURL, account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("user repo: upsert oauth account: %w", err)
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	q := `UPDATE users
	      SET full_name=$1, avatar_url=$2, provider=$3, provider_id=$4,
	          email_verified_at=$5, last_login_at=$6, updated_at=$7
	      WHERE id=$8`
	tag, err := r.db.Exec(ctx, q,
		u.FullName, u.AvatarURL, u.Provider, u.ProviderID,
		u.EmailVerifiedAt, u.LastLoginAt, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return fmt.Errorf("user repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) MarkEmailVerified(ctx context.Context, id uuid.UUID, verifiedAt time.Time) error {
	q := `UPDATE users SET email_verified_at=$1, updated_at=$1 WHERE id=$2`
	tag, err := r.db.Exec(ctx, q, verifiedAt, id)
	if err != nil {
		return fmt.Errorf("user repo: mark email verified: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string, updatedAt time.Time) error {
	q := `UPDATE users SET password_hash=$1, provider=$2, updated_at=$3 WHERE id=$4`
	tag, err := r.db.Exec(ctx, q, passwordHash, user.ProviderLocal, updatedAt, id)
	if err != nil {
		return fmt.Errorf("user repo: update password: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) RecordLogin(ctx context.Context, id uuid.UUID, loggedInAt time.Time) error {
	q := `UPDATE users SET last_login_at=$1 WHERE id=$2`
	tag, err := r.db.Exec(ctx, q, loggedInAt, id)
	if err != nil {
		return fmt.Errorf("user repo: record login: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM users WHERE id=$1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("user repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) scanOne(row pgx.Row) (*user.User, error) {
	u := &user.User{}
	err := row.Scan(
		&u.ID, &u.FullName, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Plan, &u.Provider, &u.ProviderID,
		&u.EmailVerifiedAt, &u.LastLoginAt,
		&u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("user repo: scan: %w", err)
	}
	return u, nil
}

// nullableString returns nil for an empty string so password_hash stays NULL for OAuth users.
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
