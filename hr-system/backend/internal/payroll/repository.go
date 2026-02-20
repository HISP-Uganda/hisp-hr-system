package payroll

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

func (r *Repository) ListBatches(ctx context.Context, filter BatchFilter) ([]Batch, error) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)
	position := 1

	if month := strings.TrimSpace(filter.Month); month != "" {
		conditions = append(conditions, "month = $"+fmt.Sprintf("%d", position))
		args = append(args, month)
		position++
	}
	if status := strings.TrimSpace(filter.Status); status != "" {
		conditions = append(conditions, "status = $"+fmt.Sprintf("%d", position))
		args = append(args, status)
		position++
	}

	query := `
		SELECT id, month, status, created_by, created_at, approved_by, approved_at, locked_at
		FROM payroll_batches
	`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY month DESC, id DESC"

	items := make([]Batch, 0)
	if err := r.db.SelectContext(ctx, &items, query, args...); err != nil {
		return nil, fmt.Errorf("list payroll batches: %w", err)
	}
	return items, nil
}

func (r *Repository) GetBatch(ctx context.Context, batchID int64) (Batch, error) {
	const query = `
		SELECT id, month, status, created_by, created_at, approved_by, approved_at, locked_at
		FROM payroll_batches
		WHERE id = $1
	`
	var item Batch
	if err := r.db.GetContext(ctx, &item, query, batchID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Batch{}, ErrBatchNotFound
		}
		return Batch{}, fmt.Errorf("get payroll batch: %w", err)
	}
	return item, nil
}

func (r *Repository) GetBatchEntries(ctx context.Context, batchID int64) ([]Entry, error) {
	const query = `
		SELECT
			pe.id,
			pe.batch_id,
			pe.employee_id,
			TRIM(e.last_name || ', ' || e.first_name) AS employee_name,
			pe.base_salary,
			pe.allowances_total,
			pe.deductions_total,
			pe.tax_total,
			pe.gross_pay,
			pe.net_pay,
			pe.created_at,
			pe.updated_at
		FROM payroll_entries pe
		JOIN employees e ON e.id = pe.employee_id
		WHERE pe.batch_id = $1
		ORDER BY e.last_name ASC, e.first_name ASC, pe.id ASC
	`
	items := make([]Entry, 0)
	if err := r.db.SelectContext(ctx, &items, query, batchID); err != nil {
		return nil, fmt.Errorf("list payroll entries: %w", err)
	}
	return items, nil
}

func (r *Repository) GetEntry(ctx context.Context, entryID int64) (Entry, error) {
	const query = `
		SELECT
			pe.id,
			pe.batch_id,
			pe.employee_id,
			TRIM(e.last_name || ', ' || e.first_name) AS employee_name,
			pe.base_salary,
			pe.allowances_total,
			pe.deductions_total,
			pe.tax_total,
			pe.gross_pay,
			pe.net_pay,
			pe.created_at,
			pe.updated_at
		FROM payroll_entries pe
		JOIN employees e ON e.id = pe.employee_id
		WHERE pe.id = $1
	`
	var item Entry
	if err := r.db.GetContext(ctx, &item, query, entryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrEntryNotFound
		}
		return Entry{}, fmt.Errorf("get payroll entry: %w", err)
	}
	return item, nil
}

func (r *Repository) CreateBatch(ctx context.Context, month string, createdBy int64) (Batch, error) {
	const query = `
		INSERT INTO payroll_batches (month, status, created_by)
		VALUES ($1, $2, $3)
		RETURNING id, month, status, created_by, created_at, approved_by, approved_at, locked_at
	`
	var batch Batch
	if err := r.db.GetContext(ctx, &batch, query, month, StatusDraft, createdBy); err != nil {
		if isUniqueViolation(err, "uq_payroll_batches_month") {
			return Batch{}, ErrBatchAlreadyExists
		}
		return Batch{}, fmt.Errorf("create payroll batch: %w", err)
	}
	return batch, nil
}

func (r *Repository) GenerateEntriesForBatch(ctx context.Context, batchID int64) error {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return fmt.Errorf("begin payroll generation tx: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const batchQuery = `
		SELECT id, month, status, created_by, created_at, approved_by, approved_at, locked_at
		FROM payroll_batches
		WHERE id = $1
		FOR UPDATE
	`
	var batch Batch
	if err := tx.GetContext(ctx, &batch, batchQuery, batchID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBatchNotFound
		}
		return fmt.Errorf("load payroll batch for generation: %w", err)
	}
	if batch.Status != StatusDraft {
		return ErrBatchImmutable
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM payroll_entries WHERE batch_id = $1`, batchID); err != nil {
		return fmt.Errorf("clear payroll entries for regenerate: %w", err)
	}

	type employeeBase struct {
		ID         int64   `db:"id"`
		BaseSalary float64 `db:"base_salary"`
	}
	const employeeQuery = `
		SELECT id, base_salary
		FROM employees
		WHERE LOWER(employment_status) = 'active'
		ORDER BY id ASC
	`
	employees := make([]employeeBase, 0)
	if err := tx.SelectContext(ctx, &employees, employeeQuery); err != nil {
		return fmt.Errorf("list active employees for payroll generation: %w", err)
	}

	const insertEntry = `
		INSERT INTO payroll_entries (
			batch_id,
			employee_id,
			base_salary,
			allowances_total,
			deductions_total,
			tax_total,
			gross_pay,
			net_pay
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`
	for _, employee := range employees {
		grossPay, netPay := CalculateAmounts(employee.BaseSalary, 0, 0, 0)
		if _, err := tx.ExecContext(ctx, insertEntry, batchID, employee.ID, employee.BaseSalary, 0, 0, 0, grossPay, netPay); err != nil {
			return fmt.Errorf("insert payroll entry for employee %d: %w", employee.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit payroll generation tx: %w", err)
	}
	return nil
}

func (r *Repository) UpdateEntryAmounts(ctx context.Context, entryID int64, allowancesTotal, deductionsTotal, taxTotal, grossPay, netPay float64) (Entry, error) {
	const query = `
		UPDATE payroll_entries
		SET allowances_total = $2,
			deductions_total = $3,
			tax_total = $4,
			gross_pay = $5,
			net_pay = $6,
			updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`
	var updatedID int64
	if err := r.db.GetContext(ctx, &updatedID, query, entryID, allowancesTotal, deductionsTotal, taxTotal, grossPay, netPay); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Entry{}, ErrEntryNotFound
		}
		return Entry{}, fmt.Errorf("update payroll entry amounts: %w", err)
	}
	return r.GetEntry(ctx, updatedID)
}

func (r *Repository) ApproveBatch(ctx context.Context, batchID int64, approvedBy int64, approvedAt time.Time) (Batch, error) {
	const query = `
		UPDATE payroll_batches
		SET status = $2,
			approved_by = $3,
			approved_at = $4
		WHERE id = $1
		RETURNING id, month, status, created_by, created_at, approved_by, approved_at, locked_at
	`
	var batch Batch
	if err := r.db.GetContext(ctx, &batch, query, batchID, StatusApproved, approvedBy, approvedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Batch{}, ErrBatchNotFound
		}
		return Batch{}, fmt.Errorf("approve payroll batch: %w", err)
	}
	return batch, nil
}

func (r *Repository) LockBatch(ctx context.Context, batchID int64, lockedAt time.Time) (Batch, error) {
	const query = `
		UPDATE payroll_batches
		SET status = $2,
			locked_at = $3
		WHERE id = $1
		RETURNING id, month, status, created_by, created_at, approved_by, approved_at, locked_at
	`
	var batch Batch
	if err := r.db.GetContext(ctx, &batch, query, batchID, StatusLocked, lockedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Batch{}, ErrBatchNotFound
		}
		return Batch{}, fmt.Errorf("lock payroll batch: %w", err)
	}
	return batch, nil
}

func isUniqueViolation(err error, constraint string) bool {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return false
	}
	return pqErr.Code == "23505" && pqErr.Constraint == constraint
}
