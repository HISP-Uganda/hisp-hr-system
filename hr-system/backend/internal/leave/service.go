package leave

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Actor struct {
	UserID int64
	Role   string
}

type Store interface {
	ListLeaveTypes(ctx context.Context) ([]LeaveType, error)
	CreateLeaveType(ctx context.Context, input LeaveTypeInput) (LeaveType, error)
	UpdateLeaveType(ctx context.Context, leaveTypeID int64, input LeaveTypeInput) (LeaveType, error)
	DeactivateLeaveType(ctx context.Context, leaveTypeID int64) error
	LockDate(ctx context.Context, lockDate time.Time, reason string, createdBy int64) (LockedDate, error)
	UnlockDate(ctx context.Context, lockDate time.Time) error
	ListLockedDates(ctx context.Context, year int) ([]LockedDate, error)
	ResolveEmployeeByUserID(ctx context.Context, userID int64) (int64, error)
	GetLeaveTypeByID(ctx context.Context, leaveTypeID int64) (LeaveType, error)
	GetOrCreateEntitlement(ctx context.Context, employeeID, leaveTypeID int64, year int) (LeaveEntitlement, error)
	GetUsedDays(ctx context.Context, employeeID, leaveTypeID int64, year int) (pending int, approved int, err error)
	CountApprovedOverlap(ctx context.Context, employeeID int64, startDate, endDate time.Time, excludeID *int64) (int, error)
	AnyLockedWorkingDate(ctx context.Context, days []time.Time) (bool, error)
	CreateRequest(ctx context.Context, employeeID, leaveTypeID int64, startDate, endDate time.Time, workingDays int, requestedBy int64, comment string) (LeaveRequest, error)
	GetRequestByID(ctx context.Context, requestID int64) (LeaveRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID int64, status string, actorUserID int64, comment string) (LeaveRequest, error)
	UpdateRequestByMaster(ctx context.Context, requestID int64, input ApplyInput, workingDays int) (LeaveRequest, error)
	DeleteRequest(ctx context.Context, requestID int64) error
	ListRequests(ctx context.Context, filter RequestFilter) (RequestList, error)
	ListBalances(ctx context.Context, employeeID int64, year int) ([]Balance, error)
}

type Service struct {
	store Store
}

func NewService(store Store) (*Service, error) {
	if store == nil {
		return nil, fmt.Errorf("leave store is required")
	}
	return &Service{store: store}, nil
}

func (s *Service) ListLeaveTypes(ctx context.Context) ([]LeaveType, error) {
	return s.store.ListLeaveTypes(ctx)
}

func (s *Service) CreateLeaveType(ctx context.Context, actor Actor, input LeaveTypeInput) (LeaveType, error) {
	if !isAdmin(actor.Role) {
		return LeaveType{}, ErrForbidden
	}
	if strings.TrimSpace(input.Name) == "" || input.AnnualEntitlementDays < 0 {
		return LeaveType{}, ErrInvalidInput
	}
	return s.store.CreateLeaveType(ctx, input)
}

func (s *Service) UpdateLeaveType(ctx context.Context, actor Actor, leaveTypeID int64, input LeaveTypeInput) (LeaveType, error) {
	if !isAdmin(actor.Role) {
		return LeaveType{}, ErrForbidden
	}
	if leaveTypeID <= 0 || strings.TrimSpace(input.Name) == "" || input.AnnualEntitlementDays < 0 {
		return LeaveType{}, ErrInvalidInput
	}
	return s.store.UpdateLeaveType(ctx, leaveTypeID, input)
}

func (s *Service) DeactivateLeaveType(ctx context.Context, actor Actor, leaveTypeID int64) error {
	if !isAdmin(actor.Role) {
		return ErrForbidden
	}
	if leaveTypeID <= 0 {
		return ErrInvalidInput
	}
	return s.store.DeactivateLeaveType(ctx, leaveTypeID)
}

func (s *Service) LockDate(ctx context.Context, actor Actor, input LockDateInput) (LockedDate, error) {
	if !isAdmin(actor.Role) {
		return LockedDate{}, ErrForbidden
	}
	day, err := time.Parse("2006-01-02", strings.TrimSpace(input.Date))
	if err != nil {
		return LockedDate{}, ErrInvalidInput
	}
	return s.store.LockDate(ctx, day, strings.TrimSpace(input.Reason), actor.UserID)
}

func (s *Service) UnlockDate(ctx context.Context, actor Actor, date string) error {
	if !isAdmin(actor.Role) {
		return ErrForbidden
	}
	day, err := time.Parse("2006-01-02", strings.TrimSpace(date))
	if err != nil {
		return ErrInvalidInput
	}
	return s.store.UnlockDate(ctx, day)
}

func (s *Service) ListLockedDates(ctx context.Context, actor Actor, year int) ([]LockedDate, error) {
	if year < 2000 {
		return nil, ErrInvalidInput
	}
	if !(isAdmin(actor.Role) || isStaff(actor.Role)) {
		return nil, ErrForbidden
	}
	return s.store.ListLockedDates(ctx, year)
}

func (s *Service) Apply(ctx context.Context, actor Actor, input ApplyInput) (LeaveRequest, error) {
	employeeID, err := s.resolveTargetEmployee(ctx, actor, input.EmployeeID)
	if err != nil {
		return LeaveRequest{}, err
	}
	leaveType, startDate, endDate, workingDays, workingDates, err := s.validateRequestWindow(ctx, input, nil, employeeID)
	if err != nil {
		return LeaveRequest{}, err
	}

	if leaveType.CountsTowardEntitlement {
		if err := s.validateBalance(ctx, employeeID, input.LeaveTypeID, startDate.Year(), workingDays); err != nil {
			return LeaveRequest{}, err
		}
	}
	locked, err := s.store.AnyLockedWorkingDate(ctx, workingDates)
	if err != nil {
		return LeaveRequest{}, err
	}
	if locked {
		return LeaveRequest{}, ErrLockedDate
	}

	created, err := s.store.CreateRequest(ctx, employeeID, input.LeaveTypeID, startDate, endDate, workingDays, actor.UserID, input.Comment)
	if err != nil {
		return LeaveRequest{}, err
	}
	return created, nil
}

func (s *Service) Approve(ctx context.Context, actor Actor, requestID int64, input DecisionInput) (LeaveRequest, error) {
	if !(isAdmin(actor.Role) || isHR(actor.Role)) {
		return LeaveRequest{}, ErrForbidden
	}
	request, err := s.store.GetRequestByID(ctx, requestID)
	if err != nil {
		return LeaveRequest{}, err
	}
	if request.Status != "Pending" {
		return LeaveRequest{}, ErrInvalidStatusTransition
	}

	if err := s.validateBalance(ctx, request.EmployeeID, request.LeaveTypeID, request.StartDate.Year(), request.WorkingDays); err != nil {
		return LeaveRequest{}, err
	}

	return s.store.UpdateRequestStatus(ctx, requestID, "Approved", actor.UserID, input.Comment)
}

func (s *Service) Reject(ctx context.Context, actor Actor, requestID int64, input DecisionInput) (LeaveRequest, error) {
	if !(isAdmin(actor.Role) || isHR(actor.Role)) {
		return LeaveRequest{}, ErrForbidden
	}
	request, err := s.store.GetRequestByID(ctx, requestID)
	if err != nil {
		return LeaveRequest{}, err
	}
	if request.Status != "Pending" {
		return LeaveRequest{}, ErrInvalidStatusTransition
	}
	return s.store.UpdateRequestStatus(ctx, requestID, "Rejected", actor.UserID, input.Comment)
}

func (s *Service) Cancel(ctx context.Context, actor Actor, requestID int64, input DecisionInput) (LeaveRequest, error) {
	request, err := s.store.GetRequestByID(ctx, requestID)
	if err != nil {
		return LeaveRequest{}, err
	}

	switch request.Status {
	case "Pending":
		if !(isAdmin(actor.Role) || isHR(actor.Role)) {
			selfEmployeeID, resolveErr := s.store.ResolveEmployeeByUserID(ctx, actor.UserID)
			if resolveErr != nil || selfEmployeeID != request.EmployeeID {
				return LeaveRequest{}, ErrForbidden
			}
		}
	case "Approved":
		if !(isAdmin(actor.Role) || isHR(actor.Role)) {
			return LeaveRequest{}, ErrForbidden
		}
	default:
		return LeaveRequest{}, ErrInvalidStatusTransition
	}

	return s.store.UpdateRequestStatus(ctx, requestID, "Cancelled", actor.UserID, input.Comment)
}

func (s *Service) MasterUpdate(ctx context.Context, actor Actor, requestID int64, input ApplyInput) (LeaveRequest, error) {
	if !isMaster(actor.Role) {
		return LeaveRequest{}, ErrForbidden
	}
	if requestID <= 0 {
		return LeaveRequest{}, ErrInvalidInput
	}
	if input.EmployeeID == nil || *input.EmployeeID <= 0 {
		return LeaveRequest{}, ErrInvalidInput
	}
	leaveType, startDate, endDate, workingDays, workingDates, err := s.validateRequestWindow(ctx, input, &requestID, *input.EmployeeID)
	if err != nil {
		return LeaveRequest{}, err
	}
	if leaveType.CountsTowardEntitlement {
		if err := s.validateBalance(ctx, *input.EmployeeID, input.LeaveTypeID, startDate.Year(), workingDays); err != nil {
			return LeaveRequest{}, err
		}
	}
	locked, err := s.store.AnyLockedWorkingDate(ctx, workingDates)
	if err != nil {
		return LeaveRequest{}, err
	}
	if locked {
		return LeaveRequest{}, ErrLockedDate
	}

	input.StartDate = startDate.Format("2006-01-02")
	input.EndDate = endDate.Format("2006-01-02")
	return s.store.UpdateRequestByMaster(ctx, requestID, input, workingDays)
}

func (s *Service) MasterDelete(ctx context.Context, actor Actor, requestID int64) error {
	if !isMaster(actor.Role) {
		return ErrForbidden
	}
	if requestID <= 0 {
		return ErrInvalidInput
	}
	return s.store.DeleteRequest(ctx, requestID)
}

func (s *Service) MeBalance(ctx context.Context, actor Actor, year int) (BalanceSummary, error) {
	employeeID, err := s.store.ResolveEmployeeByUserID(ctx, actor.UserID)
	if err != nil {
		return BalanceSummary{}, ErrForbidden
	}
	return s.AdminBalance(ctx, actor, employeeID, year)
}

func (s *Service) AdminBalance(ctx context.Context, actor Actor, employeeID int64, year int) (BalanceSummary, error) {
	if year < 2000 || employeeID <= 0 {
		return BalanceSummary{}, ErrInvalidInput
	}
	if !(isAdmin(actor.Role) || isHR(actor.Role) || isMaster(actor.Role)) {
		if isStaff(actor.Role) {
			selfEmployeeID, err := s.store.ResolveEmployeeByUserID(ctx, actor.UserID)
			if err != nil || selfEmployeeID != employeeID {
				return BalanceSummary{}, ErrForbidden
			}
		} else {
			return BalanceSummary{}, ErrForbidden
		}
	}
	items, err := s.store.ListBalances(ctx, employeeID, year)
	if err != nil {
		return BalanceSummary{}, err
	}
	return BalanceSummary{EmployeeID: employeeID, Year: year, Items: items}, nil
}

func (s *Service) ListRequests(ctx context.Context, actor Actor, filter RequestFilter) (RequestList, error) {
	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 100 {
		filter.PageSize = 100
	}
	if isAdmin(actor.Role) || isHR(actor.Role) || isMaster(actor.Role) {
		return s.store.ListRequests(ctx, filter)
	}
	if !isStaff(actor.Role) {
		return RequestList{}, ErrForbidden
	}
	selfEmployeeID, err := s.store.ResolveEmployeeByUserID(ctx, actor.UserID)
	if err != nil {
		return RequestList{}, ErrForbidden
	}
	filter.EmployeeID = &selfEmployeeID
	return s.store.ListRequests(ctx, filter)
}

func (s *Service) ConvertAbsenceToLeave(ctx context.Context, actor Actor, employeeID int64, absenceDate string, leaveTypeID int64) (LeaveRequest, error) {
	if !(isAdmin(actor.Role) || isHR(actor.Role)) {
		return LeaveRequest{}, ErrForbidden
	}
	// TODO: attendance module integration should verify absence source record and mark attendance state as Leave.
	return s.Apply(ctx, actor, ApplyInput{
		EmployeeID:  &employeeID,
		LeaveTypeID: leaveTypeID,
		StartDate:   absenceDate,
		EndDate:     absenceDate,
		Comment:     "absence conversion",
	})
}

func (s *Service) resolveTargetEmployee(ctx context.Context, actor Actor, requestedEmployeeID *int64) (int64, error) {
	if isAdmin(actor.Role) || isHR(actor.Role) || isMaster(actor.Role) {
		if requestedEmployeeID == nil || *requestedEmployeeID <= 0 {
			return 0, ErrInvalidInput
		}
		return *requestedEmployeeID, nil
	}
	if !isStaff(actor.Role) {
		return 0, ErrForbidden
	}
	return s.store.ResolveEmployeeByUserID(ctx, actor.UserID)
}

func (s *Service) validateRequestWindow(ctx context.Context, input ApplyInput, excludeID *int64, employeeID int64) (LeaveType, time.Time, time.Time, int, []time.Time, error) {
	if input.LeaveTypeID <= 0 {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrInvalidInput
	}
	startDate, err := time.Parse("2006-01-02", strings.TrimSpace(input.StartDate))
	if err != nil {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrInvalidInput
	}
	endDate, err := time.Parse("2006-01-02", strings.TrimSpace(input.EndDate))
	if err != nil {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrInvalidInput
	}
	if endDate.Before(startDate) {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrInvalidInput
	}

	leaveType, err := s.store.GetLeaveTypeByID(ctx, input.LeaveTypeID)
	if err != nil {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, err
	}
	if !leaveType.IsActive {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrTypeNotFound
	}

	workingDays, workingDates := ComputeWorkingDays(startDate, endDate)
	if workingDays <= 0 {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrNoWorkingDays
	}

	overlaps, err := s.store.CountApprovedOverlap(ctx, employeeID, startDate, endDate, excludeID)
	if err != nil {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, err
	}
	if overlaps > 0 {
		return LeaveType{}, time.Time{}, time.Time{}, 0, nil, ErrOverlapApproved
	}

	return leaveType, startDate, endDate, workingDays, workingDates, nil
}

func (s *Service) validateBalance(ctx context.Context, employeeID, leaveTypeID int64, year int, requestedDays int) error {
	entitlement, err := s.store.GetOrCreateEntitlement(ctx, employeeID, leaveTypeID, year)
	if err != nil {
		return err
	}
	pending, approved, err := s.store.GetUsedDays(ctx, employeeID, leaveTypeID, year)
	if err != nil {
		return err
	}
	available := CalculateAvailableBalance(entitlement.TotalDays, entitlement.ReservedDays, pending, approved)
	if requestedDays > available {
		return ErrInsufficientBalance
	}
	return nil
}

func isAdmin(role string) bool {
	return role == "Admin"
}

func isHR(role string) bool {
	return role == "HR Officer"
}

func isMaster(role string) bool {
	return role == "Master" || role == "Master Admin"
}

func isStaff(role string) bool {
	return role == "Viewer" || role == "Finance Officer"
}
