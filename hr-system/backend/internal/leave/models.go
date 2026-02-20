package leave

import "time"

type LeaveType struct {
	ID                      int64     `db:"id" json:"id"`
	Name                    string    `db:"name" json:"name"`
	AnnualEntitlementDays   int       `db:"annual_entitlement_days" json:"annual_entitlement_days"`
	IsPaid                  bool      `db:"is_paid" json:"is_paid"`
	RequiresAttachment      bool      `db:"requires_attachment" json:"requires_attachment"`
	RequiresApproval        bool      `db:"requires_approval" json:"requires_approval"`
	CountsTowardEntitlement bool      `db:"counts_toward_entitlement" json:"counts_toward_entitlement"`
	IsActive                bool      `db:"is_active" json:"is_active"`
	CreatedAt               time.Time `db:"created_at" json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at"`
}

type LeaveEntitlement struct {
	ID           int64 `db:"id"`
	EmployeeID   int64 `db:"employee_id"`
	LeaveTypeID  int64 `db:"leave_type_id"`
	Year         int   `db:"year"`
	TotalDays    int   `db:"total_days"`
	ReservedDays int   `db:"reserved_days"`
}

type LockedDate struct {
	ID        int64     `db:"id" json:"id"`
	LockDate  time.Time `db:"lock_date" json:"lock_date"`
	Reason    string    `db:"reason" json:"reason"`
	CreatedBy *int64    `db:"created_by" json:"created_by,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type LeaveRequest struct {
	ID          int64      `db:"id" json:"id"`
	EmployeeID  int64      `db:"employee_id" json:"employee_id"`
	LeaveTypeID int64      `db:"leave_type_id" json:"leave_type_id"`
	StartDate   time.Time  `db:"start_date" json:"start_date"`
	EndDate     time.Time  `db:"end_date" json:"end_date"`
	WorkingDays int        `db:"working_days" json:"working_days"`
	Status      string     `db:"status" json:"status"`
	RequestedBy *int64     `db:"requested_by" json:"requested_by,omitempty"`
	ApprovedBy  *int64     `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt  *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	RejectedBy  *int64     `db:"rejected_by" json:"rejected_by,omitempty"`
	RejectedAt  *time.Time `db:"rejected_at" json:"rejected_at,omitempty"`
	CancelledBy *int64     `db:"cancelled_by" json:"cancelled_by,omitempty"`
	CancelledAt *time.Time `db:"cancelled_at" json:"cancelled_at,omitempty"`
	Comment     string     `db:"comment" json:"comment"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`

	EmployeeName   string `db:"employee_name" json:"employee_name"`
	DepartmentID   *int64 `db:"department_id" json:"department_id,omitempty"`
	DepartmentName string `db:"department_name" json:"department_name"`
	TypeName       string `db:"type_name" json:"type_name"`
}

type LeaveTypeInput struct {
	Name                    string `json:"name"`
	AnnualEntitlementDays   int    `json:"annual_entitlement_days"`
	IsPaid                  bool   `json:"is_paid"`
	RequiresAttachment      bool   `json:"requires_attachment"`
	RequiresApproval        bool   `json:"requires_approval"`
	CountsTowardEntitlement bool   `json:"counts_toward_entitlement"`
	IsActive                bool   `json:"is_active"`
}

type ApplyInput struct {
	EmployeeID  *int64 `json:"employee_id"`
	LeaveTypeID int64  `json:"leave_type_id"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	Comment     string `json:"comment"`
}

type RequestFilter struct {
	FromDate     string `json:"from_date"`
	ToDate       string `json:"to_date"`
	DepartmentID *int64 `json:"department_id"`
	EmployeeID   *int64 `json:"employee_id"`
	LeaveTypeID  *int64 `json:"leave_type_id"`
	Status       string `json:"status"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

type RequestList struct {
	Items    []LeaveRequest `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type Balance struct {
	EmployeeID  int64   `json:"employee_id"`
	Year        int     `json:"year"`
	LeaveTypeID int64   `json:"leave_type_id"`
	TypeName    string  `json:"type_name"`
	Total       int     `json:"total"`
	Reserved    int     `json:"reserved"`
	Pending     int     `json:"pending"`
	Approved    int     `json:"approved"`
	Available   int     `json:"available"`
	UsedPercent float64 `json:"used_percent"`
}

type BalanceSummary struct {
	EmployeeID int64     `json:"employee_id"`
	Year       int       `json:"year"`
	Items      []Balance `json:"items"`
}

type DecisionInput struct {
	Comment string `json:"comment"`
}

type LockDateInput struct {
	Date   string `json:"date"`
	Reason string `json:"reason"`
}
