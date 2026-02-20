package middleware

import (
	"context"
	"strings"

	"hr-system/backend/internal/auth"
)

func RequireRoles(ctx context.Context, allowedRoles ...string) error {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return auth.ErrUnauthorized
	}

	for _, allowedRole := range allowedRoles {
		if strings.EqualFold(claims.Role, strings.TrimSpace(allowedRole)) {
			return nil
		}
	}

	return auth.ErrForbidden
}
