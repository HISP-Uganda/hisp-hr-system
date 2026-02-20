package main

import (
	"errors"
	"fmt"
	"strings"

	"hr-system/backend/bootstrap"
)

type EmployeeResponse struct {
	Success bool                   `json:"success"`
	Message string                 `json:"message"`
	Data    bootstrap.EmployeeView `json:"data"`
}

type EmployeeListResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    bootstrap.EmployeeListResult `json:"data"`
}

type DepartmentListResponse struct {
	Success bool                         `json:"success"`
	Message string                       `json:"message"`
	Data    []bootstrap.DepartmentOption `json:"data"`
}

func (a *App) CreateEmployee(accessToken string, input bootstrap.EmployeeInput) (EmployeeResponse, error) {
	if a.employees == nil || a.auth == nil {
		return EmployeeResponse{}, fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	employee, err := a.employees.CreateEmployee(a.ctx, input)
	if err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	return EmployeeResponse{
		Success: true,
		Message: "employee created",
		Data:    employee,
	}, nil
}

func (a *App) UpdateEmployee(accessToken string, employeeID int64, input bootstrap.EmployeeInput) (EmployeeResponse, error) {
	if a.employees == nil || a.auth == nil {
		return EmployeeResponse{}, fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	employee, err := a.employees.UpdateEmployee(a.ctx, employeeID, input)
	if err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	return EmployeeResponse{
		Success: true,
		Message: "employee updated",
		Data:    employee,
	}, nil
}

func (a *App) DeleteEmployee(accessToken string, employeeID int64) error {
	if a.employees == nil || a.auth == nil {
		return fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return errors.New(formatEmployeeError(err))
	}

	if err := a.employees.DeleteEmployee(a.ctx, employeeID); err != nil {
		return errors.New(formatEmployeeError(err))
	}
	return nil
}

func (a *App) GetEmployee(accessToken string, employeeID int64) (EmployeeResponse, error) {
	if a.employees == nil || a.auth == nil {
		return EmployeeResponse{}, fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	employee, err := a.employees.GetEmployee(a.ctx, employeeID)
	if err != nil {
		return EmployeeResponse{}, errors.New(formatEmployeeError(err))
	}

	return EmployeeResponse{
		Success: true,
		Message: "employee fetched",
		Data:    employee,
	}, nil
}

func (a *App) ListEmployees(accessToken string, query bootstrap.EmployeeListQuery) (EmployeeListResponse, error) {
	if a.employees == nil || a.auth == nil {
		return EmployeeListResponse{}, fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return EmployeeListResponse{}, errors.New(formatEmployeeError(err))
	}

	result, err := a.employees.ListEmployees(a.ctx, query)
	if err != nil {
		return EmployeeListResponse{}, errors.New(formatEmployeeError(err))
	}

	return EmployeeListResponse{
		Success: true,
		Message: "employees fetched",
		Data:    result,
	}, nil
}

func (a *App) ListEmployeeDepartments(accessToken string) (DepartmentListResponse, error) {
	if a.employees == nil || a.auth == nil {
		return DepartmentListResponse{}, fmt.Errorf("employee service unavailable")
	}
	if _, err := a.auth.Authorize(a.ctx, accessToken, "Admin", "HR Officer"); err != nil {
		return DepartmentListResponse{}, errors.New(formatEmployeeError(err))
	}

	result, err := a.employees.ListDepartments(a.ctx)
	if err != nil {
		return DepartmentListResponse{}, errors.New(formatEmployeeError(err))
	}

	return DepartmentListResponse{
		Success: true,
		Message: "departments fetched",
		Data:    result,
	}, nil
}

func formatEmployeeError(err error) string {
	switch {
	case bootstrap.IsUnauthorized(err):
		return "unauthorized"
	case bootstrap.IsForbidden(err):
		return "forbidden"
	case bootstrap.IsInactiveUser(err):
		return "account inactive"
	case bootstrap.IsEmployeeInvalidInput(err):
		return "invalid employee input"
	case bootstrap.IsDepartmentNotFound(err):
		return "department not found"
	case bootstrap.IsEmployeeNotFound(err):
		return "employee not found"
	default:
		return strings.TrimSpace(err.Error())
	}
}
