package main

import (
	"errors"
	"fmt"
	"strings"

	"hr-system/backend/bootstrap"
)

type UserResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Data    bootstrap.UserView `json:"data"`
}

type UserListResponse struct {
	Success bool                     `json:"success"`
	Message string                   `json:"message"`
	Data    bootstrap.UserListResult `json:"data"`
}

func (a *App) CreateUser(accessToken string, input bootstrap.UserCreateInput) (UserResponse, error) {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := a.users.CreateUser(a.ctx, actor, input)
	if err != nil {
		return UserResponse{}, errors.New(formatUserError(err))
	}

	return UserResponse{Success: true, Message: "user created", Data: user}, nil
}

func (a *App) ListUsers(accessToken string, query bootstrap.UserListQuery) (UserListResponse, error) {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return UserListResponse{}, err
	}

	result, err := a.users.ListUsers(a.ctx, actor, query)
	if err != nil {
		return UserListResponse{}, errors.New(formatUserError(err))
	}

	return UserListResponse{Success: true, Message: "users fetched", Data: result}, nil
}

func (a *App) GetUser(accessToken string, userID int64) (UserResponse, error) {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := a.users.GetUser(a.ctx, actor, userID)
	if err != nil {
		return UserResponse{}, errors.New(formatUserError(err))
	}

	return UserResponse{Success: true, Message: "user fetched", Data: user}, nil
}

func (a *App) UpdateUser(accessToken string, userID int64, input bootstrap.UserUpdateInput) (UserResponse, error) {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := a.users.UpdateUser(a.ctx, actor, userID, input)
	if err != nil {
		return UserResponse{}, errors.New(formatUserError(err))
	}

	return UserResponse{Success: true, Message: "user updated", Data: user}, nil
}

func (a *App) ResetUserPassword(accessToken string, userID int64, input bootstrap.UserResetPasswordInput) error {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return err
	}

	if err := a.users.ResetPassword(a.ctx, actor, userID, input); err != nil {
		return errors.New(formatUserError(err))
	}
	return nil
}

func (a *App) SetUserStatus(accessToken string, userID int64, input bootstrap.UserStatusInput) (UserResponse, error) {
	actor, err := a.authorizeAdminUsers(accessToken)
	if err != nil {
		return UserResponse{}, err
	}

	user, err := a.users.UpdateStatus(a.ctx, actor, userID, input)
	if err != nil {
		return UserResponse{}, errors.New(formatUserError(err))
	}

	message := "user deactivated"
	if user.IsActive {
		message = "user activated"
	}
	return UserResponse{Success: true, Message: message, Data: user}, nil
}

func (a *App) authorizeAdminUsers(accessToken string) (bootstrap.AuthUser, error) {
	if a.users == nil || a.auth == nil {
		return bootstrap.AuthUser{}, fmt.Errorf("user service unavailable")
	}
	actor, err := a.auth.Authorize(a.ctx, accessToken, "admin")
	if err != nil {
		return bootstrap.AuthUser{}, errors.New(formatUserError(err))
	}
	return actor, nil
}

func formatUserError(err error) string {
	switch {
	case bootstrap.IsUnauthorized(err), bootstrap.IsUserUnauthorized(err):
		return "401 unauthorized"
	case bootstrap.IsForbidden(err), bootstrap.IsUserForbidden(err), bootstrap.IsUserSelfDeactivateForbidden(err):
		return "403 forbidden"
	case bootstrap.IsInactiveUser(err):
		return "403 account inactive"
	case bootstrap.IsUserNotFound(err):
		return "404 not found"
	case bootstrap.IsUserConflict(err):
		return "409 conflict"
	case bootstrap.IsUserInvalidInput(err):
		return "422 invalid input"
	default:
		return strings.TrimSpace(err.Error())
	}
}
