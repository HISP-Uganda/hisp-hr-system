package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"hr-system/backend/internal/users"

	"github.com/jmoiron/sqlx"
)

type UsersFacade struct {
	service *users.Service
}

type UserCreateInput = users.CreateInput
type UserUpdateInput = users.UpdateInput
type UserResetPasswordInput = users.ResetPasswordInput
type UserStatusInput = users.StatusInput
type UserListQuery = users.ListFilter
type UserView = users.UserView
type UserListResult = users.ListResult

func NewUsersFacade(db *sqlx.DB) (*UsersFacade, error) {
	repo := users.NewRepository(db)
	service, err := users.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("create users service: %w", err)
	}
	return &UsersFacade{service: service}, nil
}

func (f *UsersFacade) CreateUser(ctx context.Context, actor AuthUser, input UserCreateInput) (UserView, error) {
	return f.service.CreateUser(ctx, toActor(actor), input)
}

func (f *UsersFacade) ListUsers(ctx context.Context, actor AuthUser, query UserListQuery) (UserListResult, error) {
	return f.service.ListUsers(ctx, toActor(actor), query)
}

func (f *UsersFacade) GetUser(ctx context.Context, actor AuthUser, userID int64) (UserView, error) {
	return f.service.GetUser(ctx, toActor(actor), userID)
}

func (f *UsersFacade) UpdateUser(ctx context.Context, actor AuthUser, userID int64, input UserUpdateInput) (UserView, error) {
	return f.service.UpdateUser(ctx, toActor(actor), userID, input)
}

func (f *UsersFacade) ResetPassword(ctx context.Context, actor AuthUser, userID int64, input UserResetPasswordInput) error {
	return f.service.ResetPassword(ctx, toActor(actor), userID, input)
}

func (f *UsersFacade) UpdateStatus(ctx context.Context, actor AuthUser, userID int64, input UserStatusInput) (UserView, error) {
	return f.service.UpdateStatus(ctx, toActor(actor), userID, input)
}

func (f *UsersFacade) SeedInitialAdmin(ctx context.Context, username, password, role string) error {
	return f.service.SeedInitialAdmin(ctx, username, password, role)
}

func IsUserUnauthorized(err error) bool {
	return errors.Is(err, users.ErrUnauthorized)
}

func IsUserForbidden(err error) bool {
	return errors.Is(err, users.ErrForbidden)
}

func IsUserInvalidInput(err error) bool {
	return errors.Is(err, users.ErrInvalidInput)
}

func IsUserNotFound(err error) bool {
	return errors.Is(err, users.ErrUserNotFound)
}

func IsUserConflict(err error) bool {
	return errors.Is(err, users.ErrUsernameExists)
}

func IsUserSelfDeactivateForbidden(err error) bool {
	return errors.Is(err, users.ErrCannotDeactivateSelf)
}

func toActor(user AuthUser) users.Actor {
	return users.Actor{
		UserID: user.ID,
		Role:   user.Role,
	}
}
