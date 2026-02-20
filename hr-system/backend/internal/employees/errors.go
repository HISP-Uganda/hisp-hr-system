package employees

import "errors"

var (
	ErrInvalidInput       = errors.New("invalid employee input")
	ErrEmployeeNotFound   = errors.New("employee not found")
	ErrDepartmentNotFound = errors.New("department not found")
)
