package leave

import "errors"

var (
	ErrInvalidInput            = errors.New("invalid leave input")
	ErrForbidden               = errors.New("forbidden")
	ErrNotFound                = errors.New("leave record not found")
	ErrTypeNotFound            = errors.New("leave type not found")
	ErrEntitlementNotFound     = errors.New("leave entitlement not found")
	ErrInsufficientBalance     = errors.New("insufficient leave balance")
	ErrLockedDate              = errors.New("requested period contains locked dates")
	ErrOverlapApproved         = errors.New("requested period overlaps approved leave")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrNoWorkingDays           = errors.New("requested period has no working days")
)
