package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"hr-system/backend/internal/payroll"

	"github.com/jmoiron/sqlx"
)

type PayrollFacade struct {
	service *payroll.Service
}

type PayrollBatch = payroll.Batch
type PayrollEntry = payroll.Entry
type PayrollBatchFilter = payroll.BatchFilter
type PayrollBatchDetail = payroll.BatchDetail
type PayrollBatchListResult = payroll.BatchListResult
type PayrollCreateBatchInput = payroll.CreateBatchInput
type PayrollUpdateEntryAmountsInput = payroll.UpdateEntryAmountsInput

func NewPayrollFacade(db *sqlx.DB) (*PayrollFacade, error) {
	repo := payroll.NewRepository(db)
	service, err := payroll.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("create payroll service: %w", err)
	}
	return &PayrollFacade{service: service}, nil
}

func (f *PayrollFacade) ListBatches(ctx context.Context, actor AuthUser, filter PayrollBatchFilter) (PayrollBatchListResult, error) {
	return f.service.ListBatches(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, filter)
}

func (f *PayrollFacade) GetBatch(ctx context.Context, actor AuthUser, batchID int64) (PayrollBatchDetail, error) {
	return f.service.GetBatch(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, batchID)
}

func (f *PayrollFacade) CreateBatch(ctx context.Context, actor AuthUser, input PayrollCreateBatchInput) (PayrollBatch, error) {
	return f.service.CreateBatch(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, input)
}

func (f *PayrollFacade) GenerateEntries(ctx context.Context, actor AuthUser, batchID int64) error {
	return f.service.GenerateEntries(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, batchID)
}

func (f *PayrollFacade) UpdateEntryAmounts(ctx context.Context, actor AuthUser, entryID int64, input PayrollUpdateEntryAmountsInput) (PayrollEntry, error) {
	return f.service.UpdateEntryAmounts(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, entryID, input)
}

func (f *PayrollFacade) ApproveBatch(ctx context.Context, actor AuthUser, batchID int64) (PayrollBatch, error) {
	return f.service.ApproveBatch(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, batchID)
}

func (f *PayrollFacade) LockBatch(ctx context.Context, actor AuthUser, batchID int64) (PayrollBatch, error) {
	return f.service.LockBatch(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, batchID)
}

func (f *PayrollFacade) ExportBatchCSV(ctx context.Context, actor AuthUser, batchID int64) (string, error) {
	return f.service.ExportBatchCSV(ctx, payroll.Actor{UserID: actor.ID, Role: actor.Role}, batchID)
}

func IsPayrollInvalidInput(err error) bool {
	return errors.Is(err, payroll.ErrInvalidInput)
}

func IsPayrollForbidden(err error) bool {
	return errors.Is(err, payroll.ErrForbidden)
}

func IsPayrollBatchNotFound(err error) bool {
	return errors.Is(err, payroll.ErrBatchNotFound)
}

func IsPayrollEntryNotFound(err error) bool {
	return errors.Is(err, payroll.ErrEntryNotFound)
}

func IsPayrollBatchExists(err error) bool {
	return errors.Is(err, payroll.ErrBatchAlreadyExists)
}

func IsPayrollStatusTransition(err error) bool {
	return errors.Is(err, payroll.ErrInvalidStatusTransition)
}

func IsPayrollImmutable(err error) bool {
	return errors.Is(err, payroll.ErrBatchImmutable)
}
