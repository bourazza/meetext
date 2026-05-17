package postgres

import (
	"context"
	"errors"
	"fmt"

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
	q := `INSERT INTO users (id, full_name, email, password_hash, avatar_url, plan, created_at, updated_at)
		  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, q,
		u.ID, u.FullName, u.Email, u.PasswordHash,
		u.AvatarURL, u.Plan, u.CreatedAt, u.UpdatedAt,
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
	q := `SELECT id, full_name, email, password_hash, avatar_url, plan, created_at, updated_at
		  FROM users WHERE id = $1`
	u := &user.User{}
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID, &u.FullName, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Plan, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("user repo: get by id: %w", err)
	}
	return u, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	q := `SELECT id, full_name, email, password_hash, avatar_url, plan, created_at, updated_at
		  FROM users WHERE email = $1`
	u := &user.User{}
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID, &u.FullName, &u.Email, &u.PasswordHash,
		&u.AvatarURL, &u.Plan, &u.CreatedAt, &u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperr.ErrNotFound
		}
		return nil, fmt.Errorf("user repo: get by email: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Update(ctx context.Context, u *user.User) error {
	q := `UPDATE users SET full_name=$1, avatar_url=$2, updated_at=$3 WHERE id=$4`
	tag, err := r.db.Exec(ctx, q, u.FullName, u.AvatarURL, u.UpdatedAt, u.ID)
	if err != nil {
		return fmt.Errorf("user repo: update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := `DELETE FROM users WHERE id = $1`
	tag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return fmt.Errorf("user repo: delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return apperr.ErrNotFound
	}
	return nil
}
