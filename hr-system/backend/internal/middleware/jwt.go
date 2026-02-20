package middleware

import (
	"context"
	"strings"

	"hr-system/backend/internal/auth"
)

type contextKey string

const (
	claimsContextKey contextKey = "auth.claims"
	userContextKey   contextKey = "auth.user"
)

type JWT struct {
	repo   *auth.Repository
	tokens *auth.TokenManager
}

func NewJWT(repo *auth.Repository, tokens *auth.TokenManager) *JWT {
	return &JWT{
		repo:   repo,
		tokens: tokens,
	}
}

// Authenticate validates JWT signature + claims, loads the user from DB, and ensures the account is active.
func (m *JWT) Authenticate(ctx context.Context, accessToken string) (context.Context, error) {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken == "" {
		return nil, auth.ErrUnauthorized
	}

	claims, err := m.tokens.ParseAccessToken(accessToken)
	if err != nil {
		return nil, auth.ErrUnauthorized
	}

	user, err := m.repo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, auth.ErrUnauthorized
	}
	if !user.IsActive {
		return nil, auth.ErrUnauthorized
	}

	authCtx := context.WithValue(ctx, claimsContextKey, claims)
	authCtx = context.WithValue(authCtx, userContextKey, auth.ToAuthUser(user))
	return authCtx, nil
}

func ClaimsFromContext(ctx context.Context) (*auth.AuthClaims, bool) {
	claims, ok := ctx.Value(claimsContextKey).(*auth.AuthClaims)
	return claims, ok
}

func UserFromContext(ctx context.Context) (auth.AuthUser, bool) {
	user, ok := ctx.Value(userContextKey).(auth.AuthUser)
	return user, ok
}
