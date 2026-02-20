package payroll

import "errors"

var (
	ErrInvalidInput            = errors.New("invalid payroll input")
	ErrForbidden               = errors.New("forbidden")
	ErrBatchNotFound           = errors.New("payroll batch not found")
	ErrEntryNotFound           = errors.New("payroll entry not found")
	ErrBatchAlreadyExists      = errors.New("payroll batch already exists")
	ErrInvalidStatusTransition = errors.New("invalid payroll status transition")
	ErrBatchImmutable          = errors.New("payroll batch is immutable")
)
