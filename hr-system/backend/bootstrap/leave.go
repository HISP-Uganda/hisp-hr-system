package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"hr-system/backend/internal/leave"

	"github.com/jmoiron/sqlx"
)

type LeaveFacade struct {
	service *leave.Service
}

type LeaveType = leave.LeaveType
type LeaveTypeInput = leave.LeaveTypeInput
type LeaveRequest = leave.LeaveRequest
type LeaveApplyInput = leave.ApplyInput
type LeaveDecisionInput = leave.DecisionInput
type LeaveRequestFilter = leave.RequestFilter
type LeaveRequestList = leave.RequestList
type LeaveBalanceSummary = leave.BalanceSummary
type LeaveLockedDate = leave.LockedDate
type LeaveLockDateInput = leave.LockDateInput

func NewLeaveFacade(db *sqlx.DB) (*LeaveFacade, error) {
	repo := leave.NewRepository(db)
	service, err := leave.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("create leave service: %w", err)
	}
	return &LeaveFacade{service: service}, nil
}

func (f *LeaveFacade) ListLeaveTypes(ctx context.Context) ([]LeaveType, error) {
	return f.service.ListLeaveTypes(ctx)
}

func (f *LeaveFacade) CreateLeaveType(ctx context.Context, actor AuthUser, input LeaveTypeInput) (LeaveType, error) {
	return f.service.CreateLeaveType(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, input)
}

func (f *LeaveFacade) UpdateLeaveType(ctx context.Context, actor AuthUser, leaveTypeID int64, input LeaveTypeInput) (LeaveType, error) {
	return f.service.UpdateLeaveType(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, leaveTypeID, input)
}

func (f *LeaveFacade) DeactivateLeaveType(ctx context.Context, actor AuthUser, leaveTypeID int64) error {
	return f.service.DeactivateLeaveType(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, leaveTypeID)
}

func (f *LeaveFacade) LockDate(ctx context.Context, actor AuthUser, input LeaveLockDateInput) (LeaveLockedDate, error) {
	return f.service.LockDate(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, input)
}

func (f *LeaveFacade) UnlockDate(ctx context.Context, actor AuthUser, date string) error {
	return f.service.UnlockDate(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, date)
}

func (f *LeaveFacade) ListLockedDates(ctx context.Context, actor AuthUser, year int) ([]LeaveLockedDate, error) {
	return f.service.ListLockedDates(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, year)
}

func (f *LeaveFacade) Apply(ctx context.Context, actor AuthUser, input LeaveApplyInput) (LeaveRequest, error) {
	return f.service.Apply(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, input)
}

func (f *LeaveFacade) Approve(ctx context.Context, actor AuthUser, requestID int64, input LeaveDecisionInput) (LeaveRequest, error) {
	return f.service.Approve(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, requestID, input)
}

func (f *LeaveFacade) Reject(ctx context.Context, actor AuthUser, requestID int64, input LeaveDecisionInput) (LeaveRequest, error) {
	return f.service.Reject(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, requestID, input)
}

func (f *LeaveFacade) Cancel(ctx context.Context, actor AuthUser, requestID int64, input LeaveDecisionInput) (LeaveRequest, error) {
	return f.service.Cancel(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, requestID, input)
}

func (f *LeaveFacade) MasterUpdate(ctx context.Context, actor AuthUser, requestID int64, input LeaveApplyInput) (LeaveRequest, error) {
	return f.service.MasterUpdate(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, requestID, input)
}

func (f *LeaveFacade) MasterDelete(ctx context.Context, actor AuthUser, requestID int64) error {
	return f.service.MasterDelete(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, requestID)
}

func (f *LeaveFacade) MeBalance(ctx context.Context, actor AuthUser, year int) (LeaveBalanceSummary, error) {
	return f.service.MeBalance(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, year)
}

func (f *LeaveFacade) AdminBalance(ctx context.Context, actor AuthUser, employeeID int64, year int) (LeaveBalanceSummary, error) {
	return f.service.AdminBalance(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, employeeID, year)
}

func (f *LeaveFacade) ListRequests(ctx context.Context, actor AuthUser, filter LeaveRequestFilter) (LeaveRequestList, error) {
	return f.service.ListRequests(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, filter)
}

func (f *LeaveFacade) ConvertAbsenceToLeave(ctx context.Context, actor AuthUser, employeeID int64, absenceDate string, leaveTypeID int64) (LeaveRequest, error) {
	return f.service.ConvertAbsenceToLeave(ctx, leave.Actor{UserID: actor.ID, Role: actor.Role}, employeeID, absenceDate, leaveTypeID)
}

func IsLeaveInvalidInput(err error) bool     { return errors.Is(err, leave.ErrInvalidInput) }
func IsLeaveForbidden(err error) bool        { return errors.Is(err, leave.ErrForbidden) }
func IsLeaveNotFound(err error) bool         { return errors.Is(err, leave.ErrNotFound) }
func IsLeaveTypeNotFound(err error) bool     { return errors.Is(err, leave.ErrTypeNotFound) }
func IsLeaveBalanceError(err error) bool     { return errors.Is(err, leave.ErrInsufficientBalance) }
func IsLeaveLockedDate(err error) bool       { return errors.Is(err, leave.ErrLockedDate) }
func IsLeaveOverlap(err error) bool          { return errors.Is(err, leave.ErrOverlapApproved) }
func IsLeaveStatusTransition(err error) bool { return errors.Is(err, leave.ErrInvalidStatusTransition) }
func IsLeaveNoWorkingDays(err error) bool    { return errors.Is(err, leave.ErrNoWorkingDays) }
