package leave

import (
	"context"
	"testing"
	"time"
)

type fakeStore struct {
	leaveType        LeaveType
	entitlement      LeaveEntitlement
	pending          int
	approved         int
	overlap          int
	locked           bool
	requestByID      LeaveRequest
	updatedStatus    string
	createdRequest   LeaveRequest
	resolvedEmployee int64
	listRequests     RequestList
}

func (f *fakeStore) ListLeaveTypes(context.Context) ([]LeaveType, error) {
	return []LeaveType{f.leaveType}, nil
}
func (f *fakeStore) CreateLeaveType(context.Context, LeaveTypeInput) (LeaveType, error) {
	return f.leaveType, nil
}
func (f *fakeStore) UpdateLeaveType(context.Context, int64, LeaveTypeInput) (LeaveType, error) {
	return f.leaveType, nil
}
func (f *fakeStore) DeactivateLeaveType(context.Context, int64) error { return nil }
func (f *fakeStore) LockDate(_ context.Context, lockDate time.Time, reason string, createdBy int64) (LockedDate, error) {
	return LockedDate{ID: 1, LockDate: lockDate, Reason: reason, CreatedBy: &createdBy}, nil
}
func (f *fakeStore) UnlockDate(context.Context, time.Time) error                { return nil }
func (f *fakeStore) ListLockedDates(context.Context, int) ([]LockedDate, error) { return nil, nil }
func (f *fakeStore) ResolveEmployeeByUserID(context.Context, int64) (int64, error) {
	if f.resolvedEmployee == 0 {
		f.resolvedEmployee = 10
	}
	return f.resolvedEmployee, nil
}
func (f *fakeStore) GetLeaveTypeByID(context.Context, int64) (LeaveType, error) {
	return f.leaveType, nil
}
func (f *fakeStore) GetOrCreateEntitlement(context.Context, int64, int64, int) (LeaveEntitlement, error) {
	return f.entitlement, nil
}
func (f *fakeStore) GetUsedDays(context.Context, int64, int64, int) (int, int, error) {
	return f.pending, f.approved, nil
}
func (f *fakeStore) CountApprovedOverlap(context.Context, int64, time.Time, time.Time, *int64) (int, error) {
	return f.overlap, nil
}
func (f *fakeStore) AnyLockedWorkingDate(context.Context, []time.Time) (bool, error) {
	return f.locked, nil
}
func (f *fakeStore) CreateRequest(context.Context, int64, int64, time.Time, time.Time, int, int64, string) (LeaveRequest, error) {
	if f.createdRequest.ID == 0 {
		f.createdRequest = LeaveRequest{ID: 77, Status: "Pending", WorkingDays: 1}
	}
	return f.createdRequest, nil
}
func (f *fakeStore) GetRequestByID(context.Context, int64) (LeaveRequest, error) {
	if f.requestByID.ID == 0 {
		f.requestByID = LeaveRequest{ID: 1, EmployeeID: 10, LeaveTypeID: 1, StartDate: time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC), Status: "Pending", WorkingDays: 1}
	}
	return f.requestByID, nil
}
func (f *fakeStore) UpdateRequestStatus(_ context.Context, _ int64, status string, _ int64, _ string) (LeaveRequest, error) {
	f.updatedStatus = status
	return LeaveRequest{ID: 1, Status: status}, nil
}
func (f *fakeStore) UpdateRequestByMaster(context.Context, int64, ApplyInput, int) (LeaveRequest, error) {
	return LeaveRequest{ID: 1, Status: "Pending"}, nil
}
func (f *fakeStore) DeleteRequest(context.Context, int64) error { return nil }
func (f *fakeStore) ListRequests(context.Context, RequestFilter) (RequestList, error) {
	return f.listRequests, nil
}
func (f *fakeStore) ListBalances(context.Context, int64, int) ([]Balance, error) {
	return []Balance{{EmployeeID: 10, Year: 2026, LeaveTypeID: 1, Total: 20, Reserved: 2, Pending: 3, Approved: 4, Available: 11}}, nil
}

func newTestService() (*Service, *fakeStore) {
	store := &fakeStore{
		leaveType:    LeaveType{ID: 1, Name: "Annual", IsActive: true, CountsTowardEntitlement: true},
		entitlement:  LeaveEntitlement{ID: 1, EmployeeID: 10, LeaveTypeID: 1, Year: 2026, TotalDays: 20, ReservedDays: 2},
		listRequests: RequestList{Items: []LeaveRequest{{ID: 1, Status: "Pending"}}, Total: 1, Page: 1, PageSize: 20},
	}
	svc, _ := NewService(store)
	return svc, store
}

func TestApplyRejectsLockedDate(t *testing.T) {
	svc, store := newTestService()
	store.locked = true

	_, err := svc.Apply(context.Background(), Actor{UserID: 1, Role: "Admin"}, ApplyInput{
		EmployeeID:  ptrInt64(10),
		LeaveTypeID: 1,
		StartDate:   "2026-02-23",
		EndDate:     "2026-02-23",
	})
	if err != ErrLockedDate {
		t.Fatalf("expected ErrLockedDate, got %v", err)
	}
}

func TestApplyRejectsInsufficientBalance(t *testing.T) {
	svc, store := newTestService()
	store.pending = 10
	store.approved = 9

	_, err := svc.Apply(context.Background(), Actor{UserID: 1, Role: "Admin"}, ApplyInput{
		EmployeeID:  ptrInt64(10),
		LeaveTypeID: 1,
		StartDate:   "2026-02-23",
		EndDate:     "2026-02-24",
	})
	if err != ErrInsufficientBalance {
		t.Fatalf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestTransitionApproveRejectCancel(t *testing.T) {
	svc, store := newTestService()

	_, err := svc.Approve(context.Background(), Actor{UserID: 2, Role: "HR Officer"}, 1, DecisionInput{})
	if err != nil || store.updatedStatus != "Approved" {
		t.Fatalf("approve failed: %v status=%s", err, store.updatedStatus)
	}

	store.requestByID = LeaveRequest{ID: 1, EmployeeID: 10, LeaveTypeID: 1, StartDate: time.Date(2026, 2, 24, 0, 0, 0, 0, time.UTC), Status: "Pending", WorkingDays: 1}
	_, err = svc.Reject(context.Background(), Actor{UserID: 2, Role: "HR Officer"}, 1, DecisionInput{})
	if err != nil || store.updatedStatus != "Rejected" {
		t.Fatalf("reject failed: %v status=%s", err, store.updatedStatus)
	}

	store.requestByID = LeaveRequest{ID: 1, EmployeeID: 10, LeaveTypeID: 1, StartDate: time.Date(2026, 2, 24, 0, 0, 0, 0, time.UTC), Status: "Pending", WorkingDays: 1}
	_, err = svc.Cancel(context.Background(), Actor{UserID: 2, Role: "HR Officer"}, 1, DecisionInput{})
	if err != nil || store.updatedStatus != "Cancelled" {
		t.Fatalf("cancel failed: %v status=%s", err, store.updatedStatus)
	}
}

func ptrInt64(v int64) *int64 { return &v }
