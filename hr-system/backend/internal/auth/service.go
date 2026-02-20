package auth

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type Service struct {
	repo   *Repository
	tokens *TokenManager
	db     *sqlx.DB
}

func NewService(db *sqlx.DB, repo *Repository, tokens *TokenManager) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("repo is required")
	}
	if tokens == nil {
		return nil, fmt.Errorf("token manager is required")
	}
	if err := tokens.ValidateSecret(); err != nil {
		return nil, err
	}
	return &Service{repo: repo, tokens: tokens, db: db}, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (AuthResult, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return AuthResult{}, ErrInvalidCredentials
	}

	user, err := s.repo.FindUserByUsername(ctx, username)
	if err != nil {
		return AuthResult{}, err
	}
	if !user.IsActive {
		return AuthResult{}, ErrInactiveUser
	}
	if !VerifyPassword(user.PasswordHash, password) {
		return AuthResult{}, ErrInvalidCredentials
	}

	accessToken, accessExpiry, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	refreshToken, refreshHash, refreshExpiry, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return AuthResult{}, err
	}
	if err := s.repo.CreateRefreshToken(ctx, user.ID, refreshHash, refreshExpiry); err != nil {
		return AuthResult{}, err
	}
	if err := s.repo.UpdateLastLoginAt(ctx, user.ID, time.Now().UTC()); err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		AccessExpiry: accessExpiry,
		User:         ToAuthUser(user),
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (AuthResult, error) {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return AuthResult{}, ErrInvalidToken
	}

	refreshHash := s.tokens.HashRefreshToken(refreshToken)
	now := time.Now().UTC()

	tx, err := s.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return AuthResult{}, fmt.Errorf("begin refresh transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	oldToken, user, err := s.repo.GetRefreshTokenAndUserForUpdate(ctx, tx, refreshHash, now)
	if err != nil {
		return AuthResult{}, err
	}
	if !user.IsActive {
		if revokeErr := s.repo.RevokeRefreshTokenByIDTx(ctx, tx, oldToken.ID, now); revokeErr != nil {
			return AuthResult{}, revokeErr
		}
		if err := tx.Commit(); err != nil {
			return AuthResult{}, fmt.Errorf("commit inactive-user revoke: %w", err)
		}
		return AuthResult{}, ErrInactiveUser
	}

	accessToken, accessExpiry, err := s.tokens.GenerateAccessToken(user)
	if err != nil {
		return AuthResult{}, err
	}
	newRefreshToken, newRefreshHash, newRefreshExpiry, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return AuthResult{}, err
	}

	if err := s.repo.RevokeRefreshTokenByIDTx(ctx, tx, oldToken.ID, now); err != nil {
		return AuthResult{}, err
	}
	if err := s.repo.CreateRefreshTokenTx(ctx, tx, user.ID, newRefreshHash, newRefreshExpiry); err != nil {
		return AuthResult{}, err
	}

	if err := tx.Commit(); err != nil {
		return AuthResult{}, fmt.Errorf("commit refresh transaction: %w", err)
	}

	return AuthResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		AccessExpiry: accessExpiry,
		User:         ToAuthUser(user),
	}, nil
}

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ErrInvalidToken
	}

	hash := s.tokens.HashRefreshToken(refreshToken)
	if err := s.repo.RevokeRefreshTokenByHash(ctx, hash, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

func (s *Service) GetUserByID(ctx context.Context, userID int64) (AuthUser, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return AuthUser{}, err
	}
	if !user.IsActive {
		return AuthUser{}, ErrInactiveUser
	}
	return ToAuthUser(user), nil
}

func (s *Service) TokenManager() *TokenManager {
	return s.tokens
}

func (s *Service) Repository() *Repository {
	return s.repo
}
