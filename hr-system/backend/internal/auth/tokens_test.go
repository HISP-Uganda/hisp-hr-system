package auth

import (
	"testing"
	"time"
)

func TestTokenManagerGenerateAndParseAccessToken(t *testing.T) {
	manager := NewTokenManager("test-secret", 15*time.Minute, 24*time.Hour)

	token, expiry, err := manager.GenerateAccessToken(User{
		ID:       10,
		Username: "tester",
		Role:     "Admin",
	})
	if err != nil {
		t.Fatalf("GenerateAccessToken() error = %v", err)
	}
	if token == "" {
		t.Fatalf("GenerateAccessToken() returned empty token")
	}
	if expiry.Before(time.Now().UTC()) {
		t.Fatalf("GenerateAccessToken() returned expired token")
	}

	claims, err := manager.ParseAccessToken(token)
	if err != nil {
		t.Fatalf("ParseAccessToken() error = %v", err)
	}
	if claims.UserID != 10 || claims.Username != "tester" || claims.Role != "Admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("my-password")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if !VerifyPassword(hash, "my-password") {
		t.Fatalf("VerifyPassword() failed for valid password")
	}
	if VerifyPassword(hash, "wrong") {
		t.Fatalf("VerifyPassword() accepted invalid password")
	}
}
