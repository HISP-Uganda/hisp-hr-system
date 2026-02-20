package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenManager struct {
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewTokenManager(secret string, accessTTL, refreshTTL time.Duration) *TokenManager {
	return &TokenManager{
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *TokenManager) GenerateAccessToken(user User) (string, time.Time, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(m.accessTTL)

	claims := AuthClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign access token: %w", err)
	}
	return signed, expiresAt, nil
}

func (m *TokenManager) GenerateRefreshToken() (rawToken, tokenHash string, expiresAt time.Time, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", time.Time{}, fmt.Errorf("generate refresh token bytes: %w", err)
	}

	rawToken = base64.RawURLEncoding.EncodeToString(buf)
	tokenHash = m.HashRefreshToken(rawToken)
	expiresAt = time.Now().UTC().Add(m.refreshTTL)
	return rawToken, tokenHash, expiresAt, nil
}

func (m *TokenManager) HashRefreshToken(refreshToken string) string {
	sum := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(sum[:])
}

func (m *TokenManager) ParseAccessToken(accessToken string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)

	token, err := parser.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}
	if token == nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	if claims.UserID == 0 || claims.Username == "" || claims.Role == "" {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt == nil {
		return nil, ErrInvalidToken
	}
	if claims.ExpiresAt.Time.Before(time.Now().UTC()) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func (m *TokenManager) ValidateSecret() error {
	if len(m.secret) == 0 {
		return errors.New("jwt secret cannot be empty")
	}
	return nil
}
