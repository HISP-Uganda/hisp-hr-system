package users

import "errors"

var (
	ErrUnauthorized         = errors.New("unauthorized")
	ErrForbidden            = errors.New("forbidden")
	ErrInvalidInput         = errors.New("invalid input")
	ErrUserNotFound         = errors.New("user not found")
	ErrUsernameExists       = errors.New("username already exists")
	ErrCannotDeactivateSelf = errors.New("cannot deactivate your own account")
)
