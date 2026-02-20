package users

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64        `db:"id" json:"id"`
	Username     string       `db:"username" json:"username"`
	PasswordHash string       `db:"password_hash" json:"-"`
	Role         string       `db:"role" json:"role"`
	IsActive     bool         `db:"is_active" json:"is_active"`
	CreatedAt    time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time    `db:"updated_at" json:"updated_at"`
	LastLoginAt  sql.NullTime `db:"last_login_at" json:"-"`
}

type UserView struct {
	ID          int64   `json:"id"`
	Username    string  `json:"username"`
	Role        string  `json:"role"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	LastLoginAt *string `json:"last_login_at"`
}

type Actor struct {
	UserID int64
	Role   string
}

type CreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateInput struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type ResetPasswordInput struct {
	NewPassword string `json:"new_password"`
}

type StatusInput struct {
	IsActive bool `json:"is_active"`
}

type ListFilter struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Q        string `json:"q"`
}

type ListResult struct {
	Items    []UserView `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

func ToUserView(user User) UserView {
	view := UserView{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.UTC().Format(time.RFC3339),
	}
	if user.LastLoginAt.Valid {
		last := user.LastLoginAt.Time.UTC().Format(time.RFC3339)
		view.LastLoginAt = &last
	}
	return view
}
