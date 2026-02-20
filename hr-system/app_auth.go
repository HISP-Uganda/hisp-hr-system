package main

import (
	"errors"
	"fmt"
	"strings"

	"hr-system/backend/bootstrap"
)

type AuthResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    bootstrap.AuthResult `json:"data"`
}

type MeResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    bootstrap.AuthUser `json:"data"`
}

func (a *App) Login(username, password string) (AuthResponse, error) {
	if a.auth == nil {
		return AuthResponse{}, fmt.Errorf("auth service unavailable")
	}

	result, err := a.auth.Login(a.ctx, username, password)
	if err != nil {
		return AuthResponse{}, errors.New(formatAuthError(err))
	}

	return AuthResponse{
		Success: true,
		Message: "login successful",
		Data:    result,
	}, nil
}

func (a *App) Refresh(refreshToken string) (AuthResponse, error) {
	if a.auth == nil {
		return AuthResponse{}, fmt.Errorf("auth service unavailable")
	}

	result, err := a.auth.Refresh(a.ctx, refreshToken)
	if err != nil {
		return AuthResponse{}, errors.New(formatAuthError(err))
	}

	return AuthResponse{
		Success: true,
		Message: "token refreshed",
		Data:    result,
	}, nil
}

func (a *App) Logout(refreshToken string) error {
	if a.auth == nil {
		return fmt.Errorf("auth service unavailable")
	}
	if err := a.auth.Logout(a.ctx, refreshToken); err != nil {
		return errors.New(formatAuthError(err))
	}
	return nil
}

func (a *App) Me(accessToken string) (MeResponse, error) {
	if a.auth == nil {
		return MeResponse{}, fmt.Errorf("auth service unavailable")
	}

	user, err := a.auth.Me(a.ctx, accessToken)
	if err != nil {
		return MeResponse{}, errors.New(formatAuthError(err))
	}

	return MeResponse{
		Success: true,
		Message: "user fetched",
		Data:    user,
	}, nil
}

func formatAuthError(err error) string {
	switch {
	case bootstrap.IsUnauthorized(err):
		return "unauthorized"
	case bootstrap.IsForbidden(err):
		return "forbidden"
	case bootstrap.IsInactiveUser(err):
		return "account inactive"
	default:
		return strings.TrimSpace(err.Error())
	}
}
