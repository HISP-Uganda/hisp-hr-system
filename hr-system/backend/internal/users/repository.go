package users

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateUser(ctx context.Context, username, passwordHash, role string) (User, error) {
	const query = `
		INSERT INTO users (username, password_hash, role, is_active)
		VALUES ($1, $2, $3, TRUE)
		RETURNING id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
	`

	var row User
	if err := r.db.GetContext(ctx, &row, query, username, passwordHash, role); err != nil {
		if isUniqueViolation(err) {
			return User{}, ErrUsernameExists
		}
		return User{}, fmt.Errorf("create user: %w", err)
	}
	return row, nil
}

func (r *Repository) GetUser(ctx context.Context, userID int64) (User, error) {
	const query = `
		SELECT id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	var row User
	if err := r.db.GetContext(ctx, &row, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("get user: %w", err)
	}
	return row, nil
}

func (r *Repository) ListUsers(ctx context.Context, filter ListFilter) ([]User, int64, error) {
	const listQuery = `
		SELECT id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%')
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	const countQuery = `
		SELECT COUNT(*)
		FROM users
		WHERE ($1 = '' OR username ILIKE '%' || $1 || '%')
	`

	offset := (filter.Page - 1) * filter.PageSize
	rows := make([]User, 0)
	if err := r.db.SelectContext(ctx, &rows, listQuery, filter.Q, filter.PageSize, offset); err != nil {
		return nil, 0, fmt.Errorf("list users: %w", err)
	}

	var total int64
	if err := r.db.GetContext(ctx, &total, countQuery, filter.Q); err != nil {
		return nil, 0, fmt.Errorf("count users: %w", err)
	}
	return rows, total, nil
}

func (r *Repository) UpdateUser(ctx context.Context, userID int64, username, role string) (User, error) {
	const query = `
		UPDATE users
		SET username = $2, role = $3, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
	`
	var row User
	if err := r.db.GetContext(ctx, &row, query, userID, username, role); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		if isUniqueViolation(err) {
			return User{}, ErrUsernameExists
		}
		return User{}, fmt.Errorf("update user: %w", err)
	}
	return row, nil
}

func (r *Repository) ResetPassword(ctx context.Context, userID int64, passwordHash string) error {
	const query = `
		UPDATE users
		SET password_hash = $2, updated_at = NOW()
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, userID, passwordHash)
	if err != nil {
		return fmt.Errorf("reset password: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("reset password rows affected: %w", err)
	}
	if affected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *Repository) UpdateStatus(ctx context.Context, userID int64, isActive bool) (User, error) {
	const query = `
		UPDATE users
		SET is_active = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
	`
	var row User
	if err := r.db.GetContext(ctx, &row, query, userID, isActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("update status: %w", err)
	}
	return row, nil
}

func (r *Repository) EnsureInitialAdmin(ctx context.Context, username, passwordHash, role string) error {
	const query = `
		INSERT INTO users (username, password_hash, role, is_active)
		VALUES ($1, $2, $3, TRUE)
		ON CONFLICT (username) DO NOTHING
	`
	if _, err := r.db.ExecContext(ctx, query, username, passwordHash, role); err != nil {
		return fmt.Errorf("ensure initial admin: %w", err)
	}
	return nil
}

func (r *Repository) WriteAuditLog(ctx context.Context, actorUserID int64, action string, entityID int64, metadata map[string]any) error {
	const query = `
		INSERT INTO audit_logs (actor_user_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, 'user', $3, $4::jsonb)
	`

	payload := []byte("{}")
	if metadata != nil {
		encoded, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("marshal audit metadata: %w", err)
		}
		payload = encoded
	}

	if _, err := r.db.ExecContext(ctx, query, actorUserID, action, fmt.Sprintf("%d", entityID), string(payload)); err != nil {
		return fmt.Errorf("write audit log: %w", err)
	}
	return nil
}

func isUniqueViolation(err error) bool {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		return strings.TrimSpace(string(pqErr.Code)) == "23505"
	}
	return false
}
