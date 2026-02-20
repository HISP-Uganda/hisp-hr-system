package users

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

type fakeStore struct {
	users      map[int64]User
	nextID     int64
	auditCalls int
}

func newFakeStore() *fakeStore {
	now := time.Now().UTC()
	return &fakeStore{
		users: map[int64]User{
			1: {ID: 1, Username: "admin", Role: "admin", IsActive: true, CreatedAt: now, UpdatedAt: now},
			2: {ID: 2, Username: "jane", Role: "viewer", IsActive: true, CreatedAt: now, UpdatedAt: now},
			3: {ID: 3, Username: "john", Role: "hr", IsActive: true, CreatedAt: now, UpdatedAt: now},
		},
		nextID: 4,
	}
}

func (f *fakeStore) CreateUser(_ context.Context, username, passwordHash, role string) (User, error) {
	for _, user := range f.users {
		if user.Username == username {
			return User{}, ErrUsernameExists
		}
	}
	now := time.Now().UTC()
	created := User{ID: f.nextID, Username: username, PasswordHash: passwordHash, Role: role, IsActive: true, CreatedAt: now, UpdatedAt: now}
	f.users[created.ID] = created
	f.nextID++
	return created, nil
}

func (f *fakeStore) GetUser(_ context.Context, userID int64) (User, error) {
	user, ok := f.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	return user, nil
}

func (f *fakeStore) ListUsers(_ context.Context, filter ListFilter) ([]User, int64, error) {
	items := make([]User, 0)
	for _, user := range f.users {
		if filter.Q == "" || contains(user.Username, filter.Q) {
			items = append(items, user)
		}
	}
	start := (filter.Page - 1) * filter.PageSize
	if start >= len(items) {
		return []User{}, int64(len(items)), nil
	}
	end := start + filter.PageSize
	if end > len(items) {
		end = len(items)
	}
	return items[start:end], int64(len(items)), nil
}

func (f *fakeStore) UpdateUser(_ context.Context, userID int64, username, role string) (User, error) {
	for id, user := range f.users {
		if id != userID && user.Username == username {
			return User{}, ErrUsernameExists
		}
	}
	user, ok := f.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	user.Username = username
	user.Role = role
	user.UpdatedAt = time.Now().UTC()
	f.users[userID] = user
	return user, nil
}

func (f *fakeStore) ResetPassword(_ context.Context, userID int64, _ string) error {
	if _, ok := f.users[userID]; !ok {
		return ErrUserNotFound
	}
	return nil
}

func (f *fakeStore) UpdateStatus(_ context.Context, userID int64, isActive bool) (User, error) {
	user, ok := f.users[userID]
	if !ok {
		return User{}, ErrUserNotFound
	}
	user.IsActive = isActive
	user.UpdatedAt = time.Now().UTC()
	f.users[userID] = user
	return user, nil
}

func (f *fakeStore) EnsureInitialAdmin(_ context.Context, _, _, _ string) error {
	return nil
}

func (f *fakeStore) WriteAuditLog(_ context.Context, _ int64, _ string, _ int64, _ map[string]any) error {
	f.auditCalls++
	return nil
}

func contains(s, sub string) bool {
	return strings.Contains(s, sub)
}

func TestAdminAccessGuards(t *testing.T) {
	store := newFakeStore()
	svc, _ := NewService(store)

	_, err := svc.ListUsers(context.Background(), Actor{}, ListFilter{})
	if !errors.Is(err, ErrUnauthorized) {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}

	_, err = svc.ListUsers(context.Background(), Actor{UserID: 2, Role: "viewer"}, ListFilter{})
	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCreateUserHappyPath(t *testing.T) {
	store := newFakeStore()
	svc, _ := NewService(store)

	user, err := svc.CreateUser(context.Background(), Actor{UserID: 1, Role: "admin"}, CreateInput{
		Username: "new-admin",
		Password: "password123",
		Role:     "admin",
	})
	if err != nil {
		t.Fatalf("create user failed: %v", err)
	}
	if user.Username != "new-admin" {
		t.Fatalf("unexpected username: %s", user.Username)
	}
	if store.auditCalls == 0 {
		t.Fatalf("expected audit log write")
	}
}

func TestDeactivateSelfForbidden(t *testing.T) {
	store := newFakeStore()
	svc, _ := NewService(store)

	_, err := svc.UpdateStatus(context.Background(), Actor{UserID: 1, Role: "admin"}, 1, StatusInput{IsActive: false})
	if !errors.Is(err, ErrCannotDeactivateSelf) {
		t.Fatalf("expected ErrCannotDeactivateSelf, got %v", err)
	}
}

func TestListPaginationAndSearch(t *testing.T) {
	store := newFakeStore()
	svc, _ := NewService(store)

	result, err := svc.ListUsers(context.Background(), Actor{UserID: 1, Role: "admin"}, ListFilter{Page: 1, PageSize: 1, Q: "ja"})
	if err != nil {
		t.Fatalf("list users failed: %v", err)
	}
	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}
	if len(result.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result.Items))
	}
	if result.Items[0].Username != "jane" {
		t.Fatalf("expected jane, got %s", result.Items[0].Username)
	}
}
