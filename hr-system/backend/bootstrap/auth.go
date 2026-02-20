package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"time"

	"hr-system/backend/internal/auth"
	"hr-system/backend/internal/middleware"

	"github.com/jmoiron/sqlx"
)

type AuthFacade struct {
	service *auth.Service
	jwt     *middleware.JWT
}

type AuthUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type AuthResult struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	AccessExpiry string   `json:"access_expiry"`
	User         AuthUser `json:"user"`
}

func NewAuthFacade(db *sqlx.DB, jwtSecret string, accessTTLSeconds, refreshTTLSeconds int64) (*AuthFacade, error) {
	repo := auth.NewRepository(db)
	tokenManager := auth.NewTokenManager(
		jwtSecret,
		secondsToDuration(accessTTLSeconds),
		secondsToDuration(refreshTTLSeconds),
	)

	service, err := auth.NewService(db, repo, tokenManager)
	if err != nil {
		return nil, fmt.Errorf("create auth service: %w", err)
	}

	return &AuthFacade{
		service: service,
		jwt:     middleware.NewJWT(repo, tokenManager),
	}, nil
}

func (f *AuthFacade) Login(ctx context.Context, username, password string) (AuthResult, error) {
	res, err := f.service.Login(ctx, username, password)
	if err != nil {
		return AuthResult{}, err
	}
	return toAuthResult(res), nil
}

func (f *AuthFacade) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	res, err := f.service.Refresh(ctx, refreshToken)
	if err != nil {
		return AuthResult{}, err
	}
	return toAuthResult(res), nil
}

func (f *AuthFacade) Logout(ctx context.Context, refreshToken string) error {
	return f.service.Logout(ctx, refreshToken)
}

func (f *AuthFacade) Me(ctx context.Context, accessToken string) (AuthUser, error) {
	authCtx, err := f.jwt.Authenticate(ctx, accessToken)
	if err != nil {
		return AuthUser{}, err
	}

	user, ok := middleware.UserFromContext(authCtx)
	if !ok {
		return AuthUser{}, auth.ErrUnauthorized
	}
	return toAuthUser(user), nil
}

func (f *AuthFacade) Authorize(ctx context.Context, accessToken string, allowedRoles ...string) (AuthUser, error) {
	authCtx, err := f.jwt.Authenticate(ctx, accessToken)
	if err != nil {
		return AuthUser{}, err
	}
	if err := middleware.RequireRoles(authCtx, allowedRoles...); err != nil {
		return AuthUser{}, err
	}

	user, ok := middleware.UserFromContext(authCtx)
	if !ok {
		return AuthUser{}, auth.ErrUnauthorized
	}
	return toAuthUser(user), nil
}

func IsUnauthorized(err error) bool {
	return errors.Is(err, auth.ErrUnauthorized) || errors.Is(err, auth.ErrInvalidToken) || errors.Is(err, auth.ErrInvalidCredentials)
}

func IsForbidden(err error) bool {
	return errors.Is(err, auth.ErrForbidden)
}

func IsInactiveUser(err error) bool {
	return errors.Is(err, auth.ErrInactiveUser)
}

func secondsToDuration(seconds int64) time.Duration {
	if seconds <= 0 {
		return 15 * time.Minute
	}
	return time.Duration(seconds) * time.Second
}

func toAuthResult(res auth.AuthResult) AuthResult {
	return AuthResult{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		AccessExpiry: res.AccessExpiry.UTC().Format("2006-01-02T15:04:05Z07:00"),
		User:         toAuthUser(res.User),
	}
}

func toAuthUser(user auth.AuthUser) AuthUser {
	return AuthUser{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}
