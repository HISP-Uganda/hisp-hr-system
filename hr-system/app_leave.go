package main

import (
	"errors"
	"fmt"
	"strings"

	"hr-system/backend/bootstrap"
)

type LeaveTypeListResponse struct {
	Success bool                  `json:"success"`
	Message string                `json:"message"`
	Data    []bootstrap.LeaveType `json:"data"`
}

type LeaveTypeResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Data    bootstrap.LeaveType `json:"data"`
}

type LeaveRequestResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    bootstrap.LeaveRequest `json:"data"`
}

type LeaveRequestListResponse struct {
	Success bool                       `json:"success"`
	Message string                     `json:"message"`
	Data    bootstrap.LeaveRequestList `json:"data"`
}

type LeaveBalanceResponse struct {
	Success bool                          `json:"success"`
	Message string                        `json:"message"`
	Data    bootstrap.LeaveBalanceSummary `json:"data"`
}

type LockedDateListResponse struct {
	Success bool                        `json:"success"`
	Message string                      `json:"message"`
	Data    []bootstrap.LeaveLockedDate `json:"data"`
}

type LockedDateResponse struct {
	Success bool                      `json:"success"`
	Message string                    `json:"message"`
	Data    bootstrap.LeaveLockedDate `json:"data"`
}

func (a *App) ListLeaveTypes(accessToken string) (LeaveTypeListResponse, error) {
	if a.leave == nil || a.auth == nil {
		return LeaveTypeListResponse{}, fmt.Errorf("leave service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer", "Finance Officer", "Viewer", "Master", "Master Admin"); err != nil {
		return LeaveTypeListResponse{}, errors.New(formatLeaveError(err))
	}
	items, err := a.leave.ListLeaveTypes(a.ctx)
	if err != nil {
		return LeaveTypeListResponse{}, errors.New(formatLeaveError(err))
	}
	return LeaveTypeListResponse{Success: true, Message: "leave types fetched", Data: items}, nil
}

func (a *App) CreateLeaveType(accessToken string, input bootstrap.LeaveTypeInput) (LeaveTypeResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveTypeResponse{}, err
	}
	item, execErr := a.leave.CreateLeaveType(a.ctx, actor, input)
	if execErr != nil {
		return LeaveTypeResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveTypeResponse{Success: true, Message: "leave type created", Data: item}, nil
}

func (a *App) UpdateLeaveType(accessToken string, leaveTypeID int64, input bootstrap.LeaveTypeInput) (LeaveTypeResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveTypeResponse{}, err
	}
	item, execErr := a.leave.UpdateLeaveType(a.ctx, actor, leaveTypeID, input)
	if execErr != nil {
		return LeaveTypeResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveTypeResponse{Success: true, Message: "leave type updated", Data: item}, nil
}

func (a *App) DeactivateLeaveType(accessToken string, leaveTypeID int64) error {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return err
	}
	if execErr := a.leave.DeactivateLeaveType(a.ctx, actor, leaveTypeID); execErr != nil {
		return errors.New(formatLeaveError(execErr))
	}
	return nil
}

func (a *App) LockLeaveDate(accessToken string, input bootstrap.LeaveLockDateInput) (LockedDateResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LockedDateResponse{}, err
	}
	item, execErr := a.leave.LockDate(a.ctx, actor, input)
	if execErr != nil {
		return LockedDateResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LockedDateResponse{Success: true, Message: "date locked", Data: item}, nil
}

func (a *App) UnlockLeaveDate(accessToken string, date string) error {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return err
	}
	if execErr := a.leave.UnlockDate(a.ctx, actor, date); execErr != nil {
		return errors.New(formatLeaveError(execErr))
	}
	return nil
}

func (a *App) ListLockedLeaveDates(accessToken string, year int) (LockedDateListResponse, error) {
	if a.leave == nil || a.auth == nil {
		return LockedDateListResponse{}, fmt.Errorf("leave service unavailable")
	}
	actor, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer", "Finance Officer", "Viewer", "Master", "Master Admin")
	if err != nil {
		return LockedDateListResponse{}, errors.New(formatLeaveError(err))
	}
	items, execErr := a.leave.ListLockedDates(a.ctx, actor, year)
	if execErr != nil {
		return LockedDateListResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LockedDateListResponse{Success: true, Message: "locked dates fetched", Data: items}, nil
}

func (a *App) ApplyLeave(accessToken string, input bootstrap.LeaveApplyInput) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.Apply(a.ctx, actor, input)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "leave applied", Data: item}, nil
}

func (a *App) ApproveLeave(accessToken string, requestID int64, input bootstrap.LeaveDecisionInput) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.Approve(a.ctx, actor, requestID, input)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "leave approved", Data: item}, nil
}

func (a *App) RejectLeave(accessToken string, requestID int64, input bootstrap.LeaveDecisionInput) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.Reject(a.ctx, actor, requestID, input)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "leave rejected", Data: item}, nil
}

func (a *App) CancelLeave(accessToken string, requestID int64, input bootstrap.LeaveDecisionInput) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.Cancel(a.ctx, actor, requestID, input)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "leave cancelled", Data: item}, nil
}

func (a *App) MasterUpdateLeave(accessToken string, requestID int64, input bootstrap.LeaveApplyInput) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.MasterUpdate(a.ctx, actor, requestID, input)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "leave updated", Data: item}, nil
}

func (a *App) MasterDeleteLeave(accessToken string, requestID int64) error {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return err
	}
	if execErr := a.leave.MasterDelete(a.ctx, actor, requestID); execErr != nil {
		return errors.New(formatLeaveError(execErr))
	}
	return nil
}

func (a *App) MeLeaveBalance(accessToken string, year int) (LeaveBalanceResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveBalanceResponse{}, err
	}
	balance, execErr := a.leave.MeBalance(a.ctx, actor, year)
	if execErr != nil {
		return LeaveBalanceResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveBalanceResponse{Success: true, Message: "balance fetched", Data: balance}, nil
}

func (a *App) AdminLeaveBalance(accessToken string, employeeID int64, year int) (LeaveBalanceResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveBalanceResponse{}, err
	}
	balance, execErr := a.leave.AdminBalance(a.ctx, actor, employeeID, year)
	if execErr != nil {
		return LeaveBalanceResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveBalanceResponse{Success: true, Message: "balance fetched", Data: balance}, nil
}

func (a *App) ListLeaveRequests(accessToken string, filter bootstrap.LeaveRequestFilter) (LeaveRequestListResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestListResponse{}, err
	}
	list, execErr := a.leave.ListRequests(a.ctx, actor, filter)
	if execErr != nil {
		return LeaveRequestListResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestListResponse{Success: true, Message: "leave requests fetched", Data: list}, nil
}

func (a *App) ConvertAbsenceToLeave(accessToken string, employeeID int64, absenceDate string, leaveTypeID int64) (LeaveRequestResponse, error) {
	actor, err := a.authorizeLeave(accessToken)
	if err != nil {
		return LeaveRequestResponse{}, err
	}
	item, execErr := a.leave.ConvertAbsenceToLeave(a.ctx, actor, employeeID, absenceDate, leaveTypeID)
	if execErr != nil {
		return LeaveRequestResponse{}, errors.New(formatLeaveError(execErr))
	}
	return LeaveRequestResponse{Success: true, Message: "absence converted", Data: item}, nil
}

func (a *App) authorizeLeave(accessToken string) (bootstrap.AuthUser, error) {
	if a.leave == nil || a.auth == nil {
		return bootstrap.AuthUser{}, fmt.Errorf("leave service unavailable")
	}
	actor, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer", "Finance Officer", "Viewer", "Master", "Master Admin")
	if err != nil {
		return bootstrap.AuthUser{}, errors.New(formatLeaveError(err))
	}
	return actor, nil
}

func formatLeaveError(err error) string {
	switch {
	case bootstrap.IsUnauthorized(err):
		return "unauthorized"
	case bootstrap.IsForbidden(err), bootstrap.IsLeaveForbidden(err):
		return "forbidden"
	case bootstrap.IsInactiveUser(err):
		return "account inactive"
	case bootstrap.IsLeaveInvalidInput(err):
		return "invalid leave input"
	case bootstrap.IsLeaveTypeNotFound(err):
		return "leave type not found"
	case bootstrap.IsLeaveNotFound(err):
		return "leave request not found"
	case bootstrap.IsLeaveNoWorkingDays(err):
		return "no working days in requested range"
	case bootstrap.IsLeaveLockedDate(err):
		return "requested period contains locked date"
	case bootstrap.IsLeaveOverlap(err):
		return "requested period overlaps approved leave"
	case bootstrap.IsLeaveBalanceError(err):
		return "insufficient leave balance"
	case bootstrap.IsLeaveStatusTransition(err):
		return "invalid leave status transition"
	default:
		return strings.TrimSpace(err.Error())
	}
}
