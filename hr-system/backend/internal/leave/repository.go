package leave

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListLeaveTypes(ctx context.Context) ([]LeaveType, error) {
	const query = `
		SELECT id, name, annual_entitlement_days, is_paid, requires_attachment, requires_approval, counts_toward_entitlement, is_active, created_at, updated_at
		FROM leave_types
		ORDER BY name ASC
	`
	items := make([]LeaveType, 0)
	if err := r.db.SelectContext(ctx, &items, query); err != nil {
		return nil, fmt.Errorf("list leave types: %w", err)
	}
	return items, nil
}

func (r *Repository) CreateLeaveType(ctx context.Context, input LeaveTypeInput) (LeaveType, error) {
	const query = `
		INSERT INTO leave_types (name, annual_entitlement_days, is_paid, requires_attachment, requires_approval, counts_toward_entitlement, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, name, annual_entitlement_days, is_paid, requires_attachment, requires_approval, counts_toward_entitlement, is_active, created_at, updated_at
	`
	var item LeaveType
	if err := r.db.GetContext(ctx, &item, query,
		strings.TrimSpace(input.Name),
		input.AnnualEntitlementDays,
		input.IsPaid,
		input.RequiresAttachment,
		input.RequiresApproval,
		input.CountsTowardEntitlement,
		input.IsActive,
	); err != nil {
		return LeaveType{}, fmt.Errorf("create leave type: %w", err)
	}
	return item, nil
}

func (r *Repository) UpdateLeaveType(ctx context.Context, leaveTypeID int64, input LeaveTypeInput) (LeaveType, error) {
	const query = `
		UPDATE leave_types
		SET name=$2, annual_entitlement_days=$3, is_paid=$4, requires_attachment=$5, requires_approval=$6, counts_toward_entitlement=$7, is_active=$8, updated_at=NOW()
		WHERE id=$1
		RETURNING id, name, annual_entitlement_days, is_paid, requires_attachment, requires_approval, counts_toward_entitlement, is_active, created_at, updated_at
	`
	var item LeaveType
	if err := r.db.GetContext(ctx, &item, query,
		leaveTypeID,
		strings.TrimSpace(input.Name),
		input.AnnualEntitlementDays,
		input.IsPaid,
		input.RequiresAttachment,
		input.RequiresApproval,
		input.CountsTowardEntitlement,
		input.IsActive,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveType{}, ErrTypeNotFound
		}
		return LeaveType{}, fmt.Errorf("update leave type: %w", err)
	}
	return item, nil
}

func (r *Repository) DeactivateLeaveType(ctx context.Context, leaveTypeID int64) error {
	const query = `UPDATE leave_types SET is_active = FALSE, updated_at = NOW() WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, leaveTypeID)
	if err != nil {
		return fmt.Errorf("deactivate leave type: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read leave type update result: %w", err)
	}
	if affected == 0 {
		return ErrTypeNotFound
	}
	return nil
}

func (r *Repository) LockDate(ctx context.Context, lockDate time.Time, reason string, createdBy int64) (LockedDate, error) {
	const query = `
		INSERT INTO leave_locked_dates (lock_date, reason, created_by)
		VALUES ($1,$2,$3)
		ON CONFLICT (lock_date)
		DO UPDATE SET reason = EXCLUDED.reason, created_by = EXCLUDED.created_by
		RETURNING id, lock_date, COALESCE(reason, '') AS reason, created_by, created_at
	`
	var item LockedDate
	if err := r.db.GetContext(ctx, &item, query, lockDate, strings.TrimSpace(reason), createdBy); err != nil {
		return LockedDate{}, fmt.Errorf("lock date: %w", err)
	}
	return item, nil
}

func (r *Repository) UnlockDate(ctx context.Context, lockDate time.Time) error {
	const query = `DELETE FROM leave_locked_dates WHERE lock_date = $1`
	_, err := r.db.ExecContext(ctx, query, lockDate)
	if err != nil {
		return fmt.Errorf("unlock date: %w", err)
	}
	return nil
}

func (r *Repository) ListLockedDates(ctx context.Context, year int) ([]LockedDate, error) {
	const query = `
		SELECT id, lock_date, COALESCE(reason, '') AS reason, created_by, created_at
		FROM leave_locked_dates
		WHERE EXTRACT(YEAR FROM lock_date) = $1
		ORDER BY lock_date ASC
	`
	items := make([]LockedDate, 0)
	if err := r.db.SelectContext(ctx, &items, query, year); err != nil {
		return nil, fmt.Errorf("list locked dates: %w", err)
	}
	return items, nil
}

func (r *Repository) ResolveEmployeeByUserID(ctx context.Context, userID int64) (int64, error) {
	const query = `SELECT id FROM employees WHERE user_id = $1`
	var employeeID int64
	if err := r.db.GetContext(ctx, &employeeID, query, userID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrForbidden
		}
		return 0, fmt.Errorf("resolve employee by user id: %w", err)
	}
	return employeeID, nil
}

func (r *Repository) GetLeaveTypeByID(ctx context.Context, leaveTypeID int64) (LeaveType, error) {
	const query = `
		SELECT id, name, annual_entitlement_days, is_paid, requires_attachment, requires_approval, counts_toward_entitlement, is_active, created_at, updated_at
		FROM leave_types
		WHERE id = $1
	`
	var item LeaveType
	if err := r.db.GetContext(ctx, &item, query, leaveTypeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveType{}, ErrTypeNotFound
		}
		return LeaveType{}, fmt.Errorf("get leave type: %w", err)
	}
	return item, nil
}

func (r *Repository) GetOrCreateEntitlement(ctx context.Context, employeeID, leaveTypeID int64, year int) (LeaveEntitlement, error) {
	const upsert = `
		INSERT INTO leave_entitlements (employee_id, leave_type_id, year, total_days, reserved_days)
		SELECT $1, $2, $3, lt.annual_entitlement_days, 0
		FROM leave_types lt
		WHERE lt.id = $2
		ON CONFLICT (employee_id, leave_type_id, year)
		DO NOTHING
	`
	if _, err := r.db.ExecContext(ctx, upsert, employeeID, leaveTypeID, year); err != nil {
		return LeaveEntitlement{}, fmt.Errorf("upsert entitlement: %w", err)
	}

	const query = `
		SELECT id, employee_id, leave_type_id, year, total_days, reserved_days
		FROM leave_entitlements
		WHERE employee_id = $1 AND leave_type_id = $2 AND year = $3
	`
	var item LeaveEntitlement
	if err := r.db.GetContext(ctx, &item, query, employeeID, leaveTypeID, year); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveEntitlement{}, ErrEntitlementNotFound
		}
		return LeaveEntitlement{}, fmt.Errorf("get entitlement: %w", err)
	}
	return item, nil
}

func (r *Repository) GetUsedDays(ctx context.Context, employeeID, leaveTypeID int64, year int) (pending int, approved int, err error) {
	const query = `
		SELECT
			COALESCE(SUM(CASE WHEN status = 'Pending' THEN working_days ELSE 0 END), 0) AS pending_days,
			COALESCE(SUM(CASE WHEN status = 'Approved' THEN working_days ELSE 0 END), 0) AS approved_days
		FROM leave_requests
		WHERE employee_id = $1
		  AND leave_type_id = $2
		  AND EXTRACT(YEAR FROM start_date) = $3
		  AND status IN ('Pending', 'Approved')
	`
	row := struct {
		Pending  int `db:"pending_days"`
		Approved int `db:"approved_days"`
	}{}
	if e := r.db.GetContext(ctx, &row, query, employeeID, leaveTypeID, year); e != nil {
		return 0, 0, fmt.Errorf("get used days: %w", e)
	}
	return row.Pending, row.Approved, nil
}

func (r *Repository) CountApprovedOverlap(ctx context.Context, employeeID int64, startDate, endDate time.Time, excludeID *int64) (int, error) {
	query := `
		SELECT COUNT(1)
		FROM leave_requests
		WHERE employee_id = $1
		  AND status = 'Approved'
		  AND start_date <= $3
		  AND end_date >= $2
	`
	args := []any{employeeID, startDate, endDate}
	if excludeID != nil {
		query += " AND id <> $4"
		args = append(args, *excludeID)
	}
	var count int
	if err := r.db.GetContext(ctx, &count, query, args...); err != nil {
		return 0, fmt.Errorf("count approved overlap: %w", err)
	}
	return count, nil
}

func (r *Repository) AnyLockedWorkingDate(ctx context.Context, days []time.Time) (bool, error) {
	if len(days) == 0 {
		return false, nil
	}
	dateValues := make([]string, 0, len(days))
	for _, day := range days {
		dateValues = append(dateValues, day.Format("2006-01-02"))
	}
	const query = `SELECT EXISTS (SELECT 1 FROM leave_locked_dates WHERE lock_date = ANY($1::date[]))`
	var exists bool
	if err := r.db.GetContext(ctx, &exists, query, pq.Array(dateValues)); err != nil {
		return false, fmt.Errorf("check locked dates: %w", err)
	}
	return exists, nil
}

func (r *Repository) CreateRequest(ctx context.Context, employeeID, leaveTypeID int64, startDate, endDate time.Time, workingDays int, requestedBy int64, comment string) (LeaveRequest, error) {
	const query = `
		INSERT INTO leave_requests (employee_id, leave_type_id, start_date, end_date, days_requested, working_days, status, requested_by, comment)
		VALUES ($1,$2,$3,$4,$5,$5,'Pending',$6,$7)
		RETURNING id, employee_id, leave_type_id, start_date, end_date, working_days, status, requested_by, approved_by, approved_at, rejected_by, rejected_at, cancelled_by, cancelled_at, COALESCE(comment, '') AS comment, created_at, updated_at
	`
	var item LeaveRequest
	if err := r.db.GetContext(ctx, &item, query, employeeID, leaveTypeID, startDate, endDate, workingDays, requestedBy, strings.TrimSpace(comment)); err != nil {
		return LeaveRequest{}, fmt.Errorf("create leave request: %w", err)
	}
	return item, nil
}

func (r *Repository) GetRequestByID(ctx context.Context, requestID int64) (LeaveRequest, error) {
	const query = `
		SELECT id, employee_id, leave_type_id, start_date, end_date, working_days, status, requested_by, approved_by, approved_at, rejected_by, rejected_at, cancelled_by, cancelled_at, COALESCE(comment, '') AS comment, created_at, updated_at
		FROM leave_requests
		WHERE id = $1
	`
	var item LeaveRequest
	if err := r.db.GetContext(ctx, &item, query, requestID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveRequest{}, ErrNotFound
		}
		return LeaveRequest{}, fmt.Errorf("get leave request: %w", err)
	}
	return item, nil
}

func (r *Repository) UpdateRequestStatus(ctx context.Context, requestID int64, status string, actorUserID int64, comment string) (LeaveRequest, error) {
	now := time.Now().UTC()
	setClause := "status = $2, updated_at = $3, comment = $5"
	switch status {
	case "Approved":
		setClause += ", approved_by = $4, approved_at = $3"
	case "Rejected":
		setClause += ", rejected_by = $4, rejected_at = $3"
	case "Cancelled":
		setClause += ", cancelled_by = $4, cancelled_at = $3"
	}
	query := `
		UPDATE leave_requests
		SET ` + setClause + `
		WHERE id = $1
		RETURNING id, employee_id, leave_type_id, start_date, end_date, working_days, status, requested_by, approved_by, approved_at, rejected_by, rejected_at, cancelled_by, cancelled_at, COALESCE(comment, '') AS comment, created_at, updated_at
	`
	var item LeaveRequest
	if err := r.db.GetContext(ctx, &item, query, requestID, status, now, actorUserID, strings.TrimSpace(comment)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveRequest{}, ErrNotFound
		}
		return LeaveRequest{}, fmt.Errorf("update leave request status: %w", err)
	}
	return item, nil
}

func (r *Repository) UpdateRequestByMaster(ctx context.Context, requestID int64, input ApplyInput, workingDays int) (LeaveRequest, error) {
	const query = `
		UPDATE leave_requests
		SET employee_id = $2, leave_type_id = $3, start_date = $4, end_date = $5, working_days = $6, days_requested = $6, comment = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, employee_id, leave_type_id, start_date, end_date, working_days, status, requested_by, approved_by, approved_at, rejected_by, rejected_at, cancelled_by, cancelled_at, COALESCE(comment, '') AS comment, created_at, updated_at
	`
	var item LeaveRequest
	if err := r.db.GetContext(ctx, &item, query, requestID, *input.EmployeeID, input.LeaveTypeID, input.StartDate, input.EndDate, workingDays, strings.TrimSpace(input.Comment)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return LeaveRequest{}, ErrNotFound
		}
		return LeaveRequest{}, fmt.Errorf("update leave request by master: %w", err)
	}
	return item, nil
}

func (r *Repository) DeleteRequest(ctx context.Context, requestID int64) error {
	const query = `DELETE FROM leave_requests WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, requestID)
	if err != nil {
		return fmt.Errorf("delete leave request: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read leave delete result: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) ListRequests(ctx context.Context, filter RequestFilter) (RequestList, error) {
	where, args := buildRequestWhere(filter)
	countQuery := `
		SELECT COUNT(1)
		FROM leave_requests lr
		JOIN employees e ON e.id = lr.employee_id
		LEFT JOIN departments d ON d.id = e.department_id
		JOIN leave_types lt ON lt.id = lr.leave_type_id
	` + where

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return RequestList{}, fmt.Errorf("count leave requests: %w", err)
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	args = append(args, pageSize, offset)
	listQuery := `
		SELECT
			lr.id,
			lr.employee_id,
			lr.leave_type_id,
			lr.start_date,
			lr.end_date,
			lr.working_days,
			lr.status,
			lr.requested_by,
			lr.approved_by,
			lr.approved_at,
			lr.rejected_by,
			lr.rejected_at,
			lr.cancelled_by,
			lr.cancelled_at,
			COALESCE(lr.comment, '') AS comment,
			lr.created_at,
			lr.updated_at,
			TRIM(e.last_name || ', ' || e.first_name) AS employee_name,
			e.department_id,
			COALESCE(d.name, '') AS department_name,
			lt.name AS type_name
		FROM leave_requests lr
		JOIN employees e ON e.id = lr.employee_id
		LEFT JOIN departments d ON d.id = e.department_id
		JOIN leave_types lt ON lt.id = lr.leave_type_id
	` + where + `
		ORDER BY lr.created_at DESC, lr.id DESC
		LIMIT $` + fmt.Sprintf("%d", len(args)-1) + ` OFFSET $` + fmt.Sprintf("%d", len(args))

	items := make([]LeaveRequest, 0)
	if err := r.db.SelectContext(ctx, &items, listQuery, args...); err != nil {
		return RequestList{}, fmt.Errorf("list leave requests: %w", err)
	}

	return RequestList{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (r *Repository) ListBalances(ctx context.Context, employeeID int64, year int) ([]Balance, error) {
	const query = `
		SELECT
			le.employee_id,
			le.year,
			le.leave_type_id,
			lt.name AS type_name,
			le.total_days AS total,
			le.reserved_days AS reserved,
			COALESCE(SUM(CASE WHEN lr.status = 'Pending' THEN lr.working_days ELSE 0 END), 0) AS pending,
			COALESCE(SUM(CASE WHEN lr.status = 'Approved' THEN lr.working_days ELSE 0 END), 0) AS approved
		FROM leave_entitlements le
		JOIN leave_types lt ON lt.id = le.leave_type_id
		LEFT JOIN leave_requests lr
			ON lr.employee_id = le.employee_id
		   AND lr.leave_type_id = le.leave_type_id
		   AND EXTRACT(YEAR FROM lr.start_date) = le.year
		   AND lr.status IN ('Pending', 'Approved')
		WHERE le.employee_id = $1 AND le.year = $2
		GROUP BY le.employee_id, le.year, le.leave_type_id, lt.name, le.total_days, le.reserved_days
		ORDER BY lt.name ASC
	`
	rows := make([]Balance, 0)
	if err := r.db.SelectContext(ctx, &rows, query, employeeID, year); err != nil {
		return nil, fmt.Errorf("list leave balances: %w", err)
	}

	for i := range rows {
		rows[i].Available = rows[i].Total - rows[i].Reserved - (rows[i].Pending + rows[i].Approved)
		used := rows[i].Pending + rows[i].Approved
		if rows[i].Total > 0 {
			rows[i].UsedPercent = float64(used) / float64(rows[i].Total) * 100
		}
	}
	return rows, nil
}

func buildRequestWhere(filter RequestFilter) (string, []any) {
	conditions := make([]string, 0, 6)
	args := make([]any, 0, 6)
	n := 1

	if strings.TrimSpace(filter.FromDate) != "" {
		conditions = append(conditions, "lr.start_date >= $"+fmt.Sprintf("%d", n))
		args = append(args, filter.FromDate)
		n++
	}
	if strings.TrimSpace(filter.ToDate) != "" {
		conditions = append(conditions, "lr.end_date <= $"+fmt.Sprintf("%d", n))
		args = append(args, filter.ToDate)
		n++
	}
	if filter.DepartmentID != nil {
		conditions = append(conditions, "e.department_id = $"+fmt.Sprintf("%d", n))
		args = append(args, *filter.DepartmentID)
		n++
	}
	if filter.EmployeeID != nil {
		conditions = append(conditions, "lr.employee_id = $"+fmt.Sprintf("%d", n))
		args = append(args, *filter.EmployeeID)
		n++
	}
	if filter.LeaveTypeID != nil {
		conditions = append(conditions, "lr.leave_type_id = $"+fmt.Sprintf("%d", n))
		args = append(args, *filter.LeaveTypeID)
		n++
	}
	if strings.TrimSpace(filter.Status) != "" {
		conditions = append(conditions, "lr.status = $"+fmt.Sprintf("%d", n))
		args = append(args, strings.TrimSpace(filter.Status))
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}
