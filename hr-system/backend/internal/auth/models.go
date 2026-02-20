package auth

import (
	"database/sql"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	ID           int64        `db:"id"`
	Username     string       `db:"username"`
	PasswordHash string       `db:"password_hash"`
	Role         string       `db:"role"`
	IsActive     bool         `db:"is_active"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at"`
	LastLoginAt  sql.NullTime `db:"last_login_at"`
}

type RefreshToken struct {
	ID        int64        `db:"id"`
	UserID    int64        `db:"user_id"`
	TokenHash string       `db:"token_hash"`
	ExpiresAt time.Time    `db:"expires_at"`
	RevokedAt sql.NullTime `db:"revoked_at"`
	CreatedAt time.Time    `db:"created_at"`
}

type AuthClaims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type AuthResult struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	AccessExpiry time.Time `json:"access_expiry"`
	User         AuthUser  `json:"user"`
}

func ToAuthUser(user User) AuthUser {
	return AuthUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}
