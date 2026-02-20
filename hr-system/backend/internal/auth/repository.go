package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindUserByUsername(ctx context.Context, username string) (User, error) {
	const query = `
		SELECT id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE username = $1
	`
	var user User
	if err := r.db.GetContext(ctx, &user, query, username); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrInvalidCredentials
		}
		return User{}, fmt.Errorf("find user by username: %w", err)
	}
	return user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, userID int64) (User, error) {
	const query = `
		SELECT id, username, password_hash, role, is_active, created_at, updated_at, last_login_at
		FROM users
		WHERE id = $1
	`
	var user User
	if err := r.db.GetContext(ctx, &user, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUnauthorized
		}
		return User{}, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

func (r *Repository) CreateRefreshToken(ctx context.Context, userID int64, tokenHash string, expiresAt time.Time) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	if _, err := r.db.ExecContext(ctx, query, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("create refresh token: %w", err)
	}
	return nil
}

func (r *Repository) UpdateLastLoginAt(ctx context.Context, userID int64, at time.Time) error {
	const query = `
		UPDATE users
		SET last_login_at = $2, updated_at = NOW()
		WHERE id = $1
	`
	if _, err := r.db.ExecContext(ctx, query, userID, at); err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}

func (r *Repository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string, revokedAt time.Time) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE token_hash = $1 AND revoked_at IS NULL
	`
	if _, err := r.db.ExecContext(ctx, query, tokenHash, revokedAt); err != nil {
		return fmt.Errorf("revoke refresh token by hash: %w", err)
	}
	return nil
}

func (r *Repository) GetRefreshTokenAndUserForUpdate(
	ctx context.Context,
	tx *sqlx.Tx,
	tokenHash string,
	now time.Time,
) (RefreshToken, User, error) {
	const query = `
		SELECT
			rt.id,
			rt.user_id,
			rt.token_hash,
			rt.expires_at,
			rt.revoked_at,
			rt.created_at,
			u.username,
			u.password_hash,
			u.role,
			u.is_active,
			u.updated_at AS user_updated_at,
			u.created_at AS user_created_at,
			u.last_login_at AS user_last_login_at
		FROM refresh_tokens rt
		JOIN users u ON u.id = rt.user_id
		WHERE rt.token_hash = $1
		  AND rt.revoked_at IS NULL
		  AND rt.expires_at > $2
		FOR UPDATE
	`

	row := struct {
		ID            int64        `db:"id"`
		UserID        int64        `db:"user_id"`
		TokenHash     string       `db:"token_hash"`
		ExpiresAt     time.Time    `db:"expires_at"`
		RevokedAt     sql.NullTime `db:"revoked_at"`
		CreatedAt     time.Time    `db:"created_at"`
		Username      string       `db:"username"`
		PasswordHash  string       `db:"password_hash"`
		Role          string       `db:"role"`
		IsActive      bool         `db:"is_active"`
		UserCreatedAt time.Time    `db:"user_created_at"`
		UserUpdatedAt time.Time    `db:"user_updated_at"`
		UserLastLogin sql.NullTime `db:"user_last_login_at"`
	}{}

	if err := tx.GetContext(ctx, &row, query, tokenHash, now); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return RefreshToken{}, User{}, ErrInvalidToken
		}
		return RefreshToken{}, User{}, fmt.Errorf("get refresh token for update: %w", err)
	}

	token := RefreshToken{
		ID:        row.ID,
		UserID:    row.UserID,
		TokenHash: row.TokenHash,
		ExpiresAt: row.ExpiresAt,
		RevokedAt: row.RevokedAt,
		CreatedAt: row.CreatedAt,
	}
	user := User{
		ID:           row.UserID,
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		Role:         row.Role,
		IsActive:     row.IsActive,
		CreatedAt:    row.UserCreatedAt,
		UpdatedAt:    row.UserUpdatedAt,
		LastLoginAt:  row.UserLastLogin,
	}

	return token, user, nil
}

func (r *Repository) RevokeRefreshTokenByIDTx(ctx context.Context, tx *sqlx.Tx, tokenID int64, revokedAt time.Time) error {
	const query = `
		UPDATE refresh_tokens
		SET revoked_at = $2
		WHERE id = $1 AND revoked_at IS NULL
	`
	res, err := tx.ExecContext(ctx, query, tokenID, revokedAt)
	if err != nil {
		return fmt.Errorf("revoke refresh token by id: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read revoke result: %w", err)
	}
	if affected == 0 {
		return ErrInvalidToken
	}
	return nil
}

func (r *Repository) CreateRefreshTokenTx(
	ctx context.Context,
	tx *sqlx.Tx,
	userID int64,
	tokenHash string,
	expiresAt time.Time,
) error {
	const query = `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
		VALUES ($1, $2, $3)
	`
	if _, err := tx.ExecContext(ctx, query, userID, tokenHash, expiresAt); err != nil {
		return fmt.Errorf("create refresh token tx: %w", err)
	}
	return nil
}
