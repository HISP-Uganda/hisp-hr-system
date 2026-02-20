package main

import (
	"errors"
	"fmt"
	"strings"

	"hr-system/backend/bootstrap"
)

type PayrollBatchListResponse struct {
	Success bool                             `json:"success"`
	Message string                           `json:"message"`
	Data    bootstrap.PayrollBatchListResult `json:"data"`
}

type PayrollBatchDetailResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    bootstrap.PayrollBatchDetail `json:"data"`
}

type PayrollBatchResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    bootstrap.PayrollBatch `json:"data"`
}

type PayrollEntryResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    bootstrap.PayrollEntry `json:"data"`
}

type PayrollCSVResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func (a *App) ListPayrollBatches(accessToken string, filter bootstrap.PayrollBatchFilter) (PayrollBatchListResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollBatchListResponse{}, err
	}
	result, execErr := a.payroll.ListBatches(a.ctx, actor, filter)
	if execErr != nil {
		return PayrollBatchListResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollBatchListResponse{Success: true, Message: "payroll batches fetched", Data: result}, nil
}

func (a *App) GetPayrollBatch(accessToken string, batchID int64) (PayrollBatchDetailResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollBatchDetailResponse{}, err
	}
	result, execErr := a.payroll.GetBatch(a.ctx, actor, batchID)
	if execErr != nil {
		return PayrollBatchDetailResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollBatchDetailResponse{Success: true, Message: "payroll batch fetched", Data: result}, nil
}

func (a *App) CreatePayrollBatch(accessToken string, input bootstrap.PayrollCreateBatchInput) (PayrollBatchResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollBatchResponse{}, err
	}
	result, execErr := a.payroll.CreateBatch(a.ctx, actor, input)
	if execErr != nil {
		return PayrollBatchResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollBatchResponse{Success: true, Message: "payroll batch created", Data: result}, nil
}

func (a *App) GeneratePayrollEntries(accessToken string, batchID int64) error {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return err
	}
	if execErr := a.payroll.GenerateEntries(a.ctx, actor, batchID); execErr != nil {
		return errors.New(formatPayrollError(execErr))
	}
	return nil
}

func (a *App) UpdatePayrollEntryAmounts(accessToken string, entryID int64, input bootstrap.PayrollUpdateEntryAmountsInput) (PayrollEntryResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollEntryResponse{}, err
	}
	result, execErr := a.payroll.UpdateEntryAmounts(a.ctx, actor, entryID, input)
	if execErr != nil {
		return PayrollEntryResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollEntryResponse{Success: true, Message: "payroll entry updated", Data: result}, nil
}

func (a *App) ApprovePayrollBatch(accessToken string, batchID int64) (PayrollBatchResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollBatchResponse{}, err
	}
	result, execErr := a.payroll.ApproveBatch(a.ctx, actor, batchID)
	if execErr != nil {
		return PayrollBatchResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollBatchResponse{Success: true, Message: "payroll batch approved", Data: result}, nil
}

func (a *App) LockPayrollBatch(accessToken string, batchID int64) (PayrollBatchResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollBatchResponse{}, err
	}
	result, execErr := a.payroll.LockBatch(a.ctx, actor, batchID)
	if execErr != nil {
		return PayrollBatchResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollBatchResponse{Success: true, Message: "payroll batch locked", Data: result}, nil
}

func (a *App) ExportPayrollBatchCSV(accessToken string, batchID int64) (PayrollCSVResponse, error) {
	actor, err := a.authorizePayroll(accessToken)
	if err != nil {
		return PayrollCSVResponse{}, err
	}
	result, execErr := a.payroll.ExportBatchCSV(a.ctx, actor, batchID)
	if execErr != nil {
		return PayrollCSVResponse{}, errors.New(formatPayrollError(execErr))
	}
	return PayrollCSVResponse{Success: true, Message: "payroll csv exported", Data: result}, nil
}

func (a *App) authorizePayroll(accessToken string) (bootstrap.AuthUser, error) {
	if a.payroll == nil || a.auth == nil {
		return bootstrap.AuthUser{}, fmt.Errorf("payroll service unavailable")
	}
	actor, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "Finance Officer")
	if err != nil {
		return bootstrap.AuthUser{}, errors.New(formatPayrollError(err))
	}
	return actor, nil
}

func formatPayrollError(err error) string {
	switch {
	case bootstrap.IsUnauthorized(err):
		return "unauthorized"
	case bootstrap.IsForbidden(err), bootstrap.IsPayrollForbidden(err):
		return "forbidden"
	case bootstrap.IsInactiveUser(err):
		return "account inactive"
	case bootstrap.IsPayrollInvalidInput(err):
		return "invalid payroll input"
	case bootstrap.IsPayrollBatchNotFound(err):
		return "payroll batch not found"
	case bootstrap.IsPayrollEntryNotFound(err):
		return "payroll entry not found"
	case bootstrap.IsPayrollBatchExists(err):
		return "payroll batch already exists for month"
	case bootstrap.IsPayrollStatusTransition(err):
		return "invalid payroll status transition"
	case bootstrap.IsPayrollImmutable(err):
		return "payroll batch is immutable"
	default:
		return strings.TrimSpace(err.Error())
	}
}
