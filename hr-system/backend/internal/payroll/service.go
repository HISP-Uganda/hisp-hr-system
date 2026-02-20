package payroll

import (
	"context"
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Store interface {
	ListBatches(ctx context.Context, filter BatchFilter) ([]Batch, error)
	GetBatch(ctx context.Context, batchID int64) (Batch, error)
	GetBatchEntries(ctx context.Context, batchID int64) ([]Entry, error)
	GetEntry(ctx context.Context, entryID int64) (Entry, error)
	CreateBatch(ctx context.Context, month string, createdBy int64) (Batch, error)
	GenerateEntriesForBatch(ctx context.Context, batchID int64) error
	UpdateEntryAmounts(ctx context.Context, entryID int64, allowancesTotal, deductionsTotal, taxTotal, grossPay, netPay float64) (Entry, error)
	ApproveBatch(ctx context.Context, batchID int64, approvedBy int64, approvedAt time.Time) (Batch, error)
	LockBatch(ctx context.Context, batchID int64, lockedAt time.Time) (Batch, error)
}

type Service struct {
	store Store
}

func NewService(store Store) (*Service, error) {
	if store == nil {
		return nil, fmt.Errorf("payroll store is required")
	}
	return &Service{store: store}, nil
}

func (s *Service) ListBatches(ctx context.Context, actor Actor, filter BatchFilter) (BatchListResult, error) {
	if !canManagePayroll(actor.Role) {
		return BatchListResult{}, ErrForbidden
	}
	if filter.Month != "" && !isValidMonth(filter.Month) {
		return BatchListResult{}, ErrInvalidInput
	}
	if filter.Status != "" && !isValidStatus(filter.Status) {
		return BatchListResult{}, ErrInvalidInput
	}

	items, err := s.store.ListBatches(ctx, filter)
	if err != nil {
		return BatchListResult{}, err
	}
	return BatchListResult{Items: items}, nil
}

func (s *Service) GetBatch(ctx context.Context, actor Actor, batchID int64) (BatchDetail, error) {
	if !canManagePayroll(actor.Role) {
		return BatchDetail{}, ErrForbidden
	}
	if batchID <= 0 {
		return BatchDetail{}, ErrInvalidInput
	}

	batch, err := s.store.GetBatch(ctx, batchID)
	if err != nil {
		return BatchDetail{}, err
	}
	entries, err := s.store.GetBatchEntries(ctx, batchID)
	if err != nil {
		return BatchDetail{}, err
	}
	return BatchDetail{Batch: batch, Entries: entries}, nil
}

func (s *Service) CreateBatch(ctx context.Context, actor Actor, input CreateBatchInput) (Batch, error) {
	if !canManagePayroll(actor.Role) {
		return Batch{}, ErrForbidden
	}
	month := strings.TrimSpace(input.Month)
	if !isValidMonth(month) {
		return Batch{}, ErrInvalidInput
	}
	return s.store.CreateBatch(ctx, month, actor.UserID)
}

func (s *Service) GenerateEntries(ctx context.Context, actor Actor, batchID int64) error {
	if !canManagePayroll(actor.Role) {
		return ErrForbidden
	}
	if batchID <= 0 {
		return ErrInvalidInput
	}
	return s.store.GenerateEntriesForBatch(ctx, batchID)
}

func (s *Service) UpdateEntryAmounts(ctx context.Context, actor Actor, entryID int64, input UpdateEntryAmountsInput) (Entry, error) {
	if !canManagePayroll(actor.Role) {
		return Entry{}, ErrForbidden
	}
	if entryID <= 0 {
		return Entry{}, ErrInvalidInput
	}
	if input.AllowancesTotal < 0 || input.DeductionsTotal < 0 || input.TaxTotal < 0 {
		return Entry{}, ErrInvalidInput
	}

	entry, err := s.store.GetEntry(ctx, entryID)
	if err != nil {
		return Entry{}, err
	}
	batch, err := s.store.GetBatch(ctx, entry.BatchID)
	if err != nil {
		return Entry{}, err
	}
	if batch.Status != StatusDraft {
		return Entry{}, ErrBatchImmutable
	}

	grossPay, netPay := CalculateAmounts(entry.BaseSalary, input.AllowancesTotal, input.DeductionsTotal, input.TaxTotal)
	return s.store.UpdateEntryAmounts(ctx, entryID, input.AllowancesTotal, input.DeductionsTotal, input.TaxTotal, grossPay, netPay)
}

func (s *Service) ApproveBatch(ctx context.Context, actor Actor, batchID int64) (Batch, error) {
	if !canManagePayroll(actor.Role) {
		return Batch{}, ErrForbidden
	}
	if batchID <= 0 {
		return Batch{}, ErrInvalidInput
	}

	batch, err := s.store.GetBatch(ctx, batchID)
	if err != nil {
		return Batch{}, err
	}
	if batch.Status != StatusDraft {
		return Batch{}, ErrInvalidStatusTransition
	}
	return s.store.ApproveBatch(ctx, batchID, actor.UserID, time.Now().UTC())
}

func (s *Service) LockBatch(ctx context.Context, actor Actor, batchID int64) (Batch, error) {
	if !canManagePayroll(actor.Role) {
		return Batch{}, ErrForbidden
	}
	if batchID <= 0 {
		return Batch{}, ErrInvalidInput
	}

	batch, err := s.store.GetBatch(ctx, batchID)
	if err != nil {
		return Batch{}, err
	}
	if batch.Status != StatusApproved {
		return Batch{}, ErrInvalidStatusTransition
	}
	return s.store.LockBatch(ctx, batchID, time.Now().UTC())
}

func (s *Service) ExportBatchCSV(ctx context.Context, actor Actor, batchID int64) (string, error) {
	if !canManagePayroll(actor.Role) {
		return "", ErrForbidden
	}
	if batchID <= 0 {
		return "", ErrInvalidInput
	}

	batch, err := s.store.GetBatch(ctx, batchID)
	if err != nil {
		return "", err
	}
	if batch.Status != StatusApproved && batch.Status != StatusLocked {
		return "", ErrBatchImmutable
	}

	entries, err := s.store.GetBatchEntries(ctx, batchID)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	writer := csv.NewWriter(&sb)
	if writeErr := writer.Write([]string{"Employee Name", "Base Salary", "Allowances", "Deductions", "Tax", "Gross Pay", "Net Pay"}); writeErr != nil {
		return "", fmt.Errorf("write payroll csv header: %w", writeErr)
	}
	for _, entry := range entries {
		record := []string{
			entry.EmployeeName,
			toMoney(entry.BaseSalary),
			toMoney(entry.AllowancesTotal),
			toMoney(entry.DeductionsTotal),
			toMoney(entry.TaxTotal),
			toMoney(entry.GrossPay),
			toMoney(entry.NetPay),
		}
		if writeErr := writer.Write(record); writeErr != nil {
			return "", fmt.Errorf("write payroll csv row: %w", writeErr)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("flush payroll csv: %w", err)
	}
	return sb.String(), nil
}

func canManagePayroll(role string) bool {
	return role == "Admin" || role == "Finance Officer"
}

func isValidMonth(month string) bool {
	_, err := time.Parse("2006-01", strings.TrimSpace(month))
	return err == nil
}

func isValidStatus(status string) bool {
	switch strings.TrimSpace(status) {
	case StatusDraft, StatusApproved, StatusLocked:
		return true
	default:
		return false
	}
}

func toMoney(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}
