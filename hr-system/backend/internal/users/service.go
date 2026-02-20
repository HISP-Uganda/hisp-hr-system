package users

import (
	"context"
	"fmt"
	"strings"

	"hr-system/backend/internal/auth"
)

const minPasswordLength = 8

type Store interface {
	CreateUser(ctx context.Context, username, passwordHash, role string) (User, error)
	GetUser(ctx context.Context, userID int64) (User, error)
	ListUsers(ctx context.Context, filter ListFilter) ([]User, int64, error)
	UpdateUser(ctx context.Context, userID int64, username, role string) (User, error)
	ResetPassword(ctx context.Context, userID int64, passwordHash string) error
	UpdateStatus(ctx context.Context, userID int64, isActive bool) (User, error)
	EnsureInitialAdmin(ctx context.Context, username, passwordHash, role string) error
	WriteAuditLog(ctx context.Context, actorUserID int64, action string, entityID int64, metadata map[string]any) error
}

type Service struct {
	store Store
}

func NewService(store Store) (*Service, error) {
	if store == nil {
		return nil, fmt.Errorf("users store is required")
	}
	return &Service{store: store}, nil
}

func (s *Service) CreateUser(ctx context.Context, actor Actor, input CreateInput) (UserView, error) {
	if err := requireAdmin(actor); err != nil {
		return UserView{}, err
	}

	username, role, password, err := normalizeCreateInput(input)
	if err != nil {
		return UserView{}, err
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return UserView{}, fmt.Errorf("hash password: %w", err)
	}

	created, err := s.store.CreateUser(ctx, username, hash, role)
	if err != nil {
		return UserView{}, err
	}

	_ = s.store.WriteAuditLog(ctx, actor.UserID, "user.create", created.ID, map[string]any{
		"username": created.Username,
		"role":     created.Role,
	})

	return ToUserView(created), nil
}

func (s *Service) GetUser(ctx context.Context, actor Actor, userID int64) (UserView, error) {
	if err := requireAdmin(actor); err != nil {
		return UserView{}, err
	}
	if userID <= 0 {
		return UserView{}, ErrInvalidInput
	}

	row, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return UserView{}, err
	}
	return ToUserView(row), nil
}

func (s *Service) ListUsers(ctx context.Context, actor Actor, filter ListFilter) (ListResult, error) {
	if err := requireAdmin(actor); err != nil {
		return ListResult{}, err
	}

	normalized := filter
	if normalized.Page <= 0 {
		normalized.Page = 1
	}
	if normalized.PageSize <= 0 {
		normalized.PageSize = 10
	}
	if normalized.PageSize > 100 {
		normalized.PageSize = 100
	}
	normalized.Q = strings.TrimSpace(normalized.Q)

	rows, total, err := s.store.ListUsers(ctx, normalized)
	if err != nil {
		return ListResult{}, err
	}

	items := make([]UserView, 0, len(rows))
	for _, row := range rows {
		items = append(items, ToUserView(row))
	}

	return ListResult{
		Items:    items,
		Total:    total,
		Page:     normalized.Page,
		PageSize: normalized.PageSize,
	}, nil
}

func (s *Service) UpdateUser(ctx context.Context, actor Actor, userID int64, input UpdateInput) (UserView, error) {
	if err := requireAdmin(actor); err != nil {
		return UserView{}, err
	}
	if userID <= 0 {
		return UserView{}, ErrInvalidInput
	}

	username, role, err := normalizeUpdateInput(input)
	if err != nil {
		return UserView{}, err
	}

	updated, err := s.store.UpdateUser(ctx, userID, username, role)
	if err != nil {
		return UserView{}, err
	}

	_ = s.store.WriteAuditLog(ctx, actor.UserID, "user.update", updated.ID, map[string]any{
		"username": updated.Username,
		"role":     updated.Role,
	})

	return ToUserView(updated), nil
}

func (s *Service) ResetPassword(ctx context.Context, actor Actor, userID int64, input ResetPasswordInput) error {
	if err := requireAdmin(actor); err != nil {
		return err
	}
	if userID <= 0 {
		return ErrInvalidInput
	}

	password := strings.TrimSpace(input.NewPassword)
	if len(password) < minPasswordLength {
		return ErrInvalidInput
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := s.store.ResetPassword(ctx, userID, hash); err != nil {
		return err
	}

	_ = s.store.WriteAuditLog(ctx, actor.UserID, "user.reset_password", userID, nil)
	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, actor Actor, userID int64, input StatusInput) (UserView, error) {
	if err := requireAdmin(actor); err != nil {
		return UserView{}, err
	}
	if userID <= 0 {
		return UserView{}, ErrInvalidInput
	}
	if actor.UserID == userID && !input.IsActive {
		return UserView{}, ErrCannotDeactivateSelf
	}

	updated, err := s.store.UpdateStatus(ctx, userID, input.IsActive)
	if err != nil {
		return UserView{}, err
	}
	action := "user.deactivate"
	if updated.IsActive {
		action = "user.activate"
	}
	_ = s.store.WriteAuditLog(ctx, actor.UserID, action, updated.ID, map[string]any{
		"is_active": updated.IsActive,
	})

	return ToUserView(updated), nil
}

func (s *Service) SeedInitialAdmin(ctx context.Context, username, password, role string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	role = strings.TrimSpace(role)
	if username == "" && password == "" {
		return nil
	}
	if username == "" || password == "" {
		return ErrInvalidInput
	}
	if role == "" {
		role = "admin"
	}
	if len(password) < minPasswordLength {
		return ErrInvalidInput
	}

	hash, err := auth.HashPassword(password)
	if err != nil {
		return fmt.Errorf("hash seed admin password: %w", err)
	}
	if err := s.store.EnsureInitialAdmin(ctx, username, hash, role); err != nil {
		return err
	}
	return nil
}

func normalizeCreateInput(input CreateInput) (string, string, string, error) {
	username := strings.TrimSpace(input.Username)
	role := strings.TrimSpace(input.Role)
	password := strings.TrimSpace(input.Password)
	if username == "" || role == "" || len(password) < minPasswordLength {
		return "", "", "", ErrInvalidInput
	}
	return username, role, password, nil
}

func normalizeUpdateInput(input UpdateInput) (string, string, error) {
	username := strings.TrimSpace(input.Username)
	role := strings.TrimSpace(input.Role)
	if username == "" || role == "" {
		return "", "", ErrInvalidInput
	}
	return username, role, nil
}

func requireAdmin(actor Actor) error {
	if actor.UserID <= 0 {
		return ErrUnauthorized
	}
	if !strings.EqualFold(strings.TrimSpace(actor.Role), "admin") {
		return ErrForbidden
	}
	return nil
}
