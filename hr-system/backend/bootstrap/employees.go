package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"hr-system/backend/internal/employees"

	"github.com/jmoiron/sqlx"
)

type EmployeesFacade struct {
	service *employees.Service
}

type EmployeeInput = employees.UpsertEmployeeInput
type EmployeeListQuery = employees.EmployeeListFilter
type EmployeeView = employees.EmployeeView
type EmployeeListResult = employees.EmployeeListResult
type DepartmentOption = employees.DepartmentOption

func NewEmployeesFacade(db *sqlx.DB) (*EmployeesFacade, error) {
	repo := employees.NewRepository(db)
	service, err := employees.NewService(repo)
	if err != nil {
		return nil, fmt.Errorf("create employees service: %w", err)
	}
	return &EmployeesFacade{service: service}, nil
}

func (f *EmployeesFacade) CreateEmployee(ctx context.Context, input EmployeeInput) (EmployeeView, error) {
	return f.service.CreateEmployee(ctx, input)
}

func (f *EmployeesFacade) UpdateEmployee(ctx context.Context, employeeID int64, input EmployeeInput) (EmployeeView, error) {
	return f.service.UpdateEmployee(ctx, employeeID, input)
}

func (f *EmployeesFacade) DeleteEmployee(ctx context.Context, employeeID int64) error {
	return f.service.DeleteEmployee(ctx, employeeID)
}

func (f *EmployeesFacade) GetEmployee(ctx context.Context, employeeID int64) (EmployeeView, error) {
	return f.service.GetEmployee(ctx, employeeID)
}

func (f *EmployeesFacade) ListEmployees(ctx context.Context, query EmployeeListQuery) (EmployeeListResult, error) {
	return f.service.ListEmployees(ctx, query)
}

func (f *EmployeesFacade) ListDepartments(ctx context.Context) ([]DepartmentOption, error) {
	return f.service.ListDepartments(ctx)
}

func IsEmployeeInvalidInput(err error) bool {
	return errors.Is(err, employees.ErrInvalidInput)
}

func IsEmployeeNotFound(err error) bool {
	return errors.Is(err, employees.ErrEmployeeNotFound)
}

func IsDepartmentNotFound(err error) bool {
	return errors.Is(err, employees.ErrDepartmentNotFound)
}
