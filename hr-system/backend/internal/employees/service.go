package employees

import (
	"context"
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) (*Service, error) {
	if repo == nil {
		return nil, fmt.Errorf("employees repository is required")
	}
	return &Service{repo: repo}, nil
}

func (s *Service) CreateEmployee(ctx context.Context, input UpsertEmployeeInput) (EmployeeView, error) {
	normalized, err := s.normalizeInput(input)
	if err != nil {
		return EmployeeView{}, err
	}
	if err := s.ensureDepartmentIntegrity(ctx, normalized.DepartmentID); err != nil {
		return EmployeeView{}, err
	}

	row, err := s.repo.CreateEmployee(ctx, normalized)
	if err != nil {
		return EmployeeView{}, err
	}

	created, err := s.repo.GetEmployee(ctx, row.ID)
	if err != nil {
		return EmployeeView{}, err
	}
	return ToEmployeeView(created), nil
}

func (s *Service) UpdateEmployee(ctx context.Context, employeeID int64, input UpsertEmployeeInput) (EmployeeView, error) {
	if employeeID <= 0 {
		return EmployeeView{}, ErrInvalidInput
	}

	normalized, err := s.normalizeInput(input)
	if err != nil {
		return EmployeeView{}, err
	}
	if err := s.ensureDepartmentIntegrity(ctx, normalized.DepartmentID); err != nil {
		return EmployeeView{}, err
	}

	_, err = s.repo.UpdateEmployee(ctx, employeeID, normalized)
	if err != nil {
		return EmployeeView{}, err
	}

	updated, err := s.repo.GetEmployee(ctx, employeeID)
	if err != nil {
		return EmployeeView{}, err
	}
	return ToEmployeeView(updated), nil
}

func (s *Service) DeleteEmployee(ctx context.Context, employeeID int64) error {
	if employeeID <= 0 {
		return ErrInvalidInput
	}
	return s.repo.DeleteEmployee(ctx, employeeID)
}

func (s *Service) GetEmployee(ctx context.Context, employeeID int64) (EmployeeView, error) {
	if employeeID <= 0 {
		return EmployeeView{}, ErrInvalidInput
	}
	row, err := s.repo.GetEmployee(ctx, employeeID)
	if err != nil {
		return EmployeeView{}, err
	}
	return ToEmployeeView(row), nil
}

func (s *Service) ListEmployees(ctx context.Context, filter EmployeeListFilter) (EmployeeListResult, error) {
	normalizedFilter := filter
	if normalizedFilter.Page <= 0 {
		normalizedFilter.Page = 1
	}
	if normalizedFilter.PageSize <= 0 {
		normalizedFilter.PageSize = 10
	}
	if normalizedFilter.PageSize > 100 {
		normalizedFilter.PageSize = 100
	}
	if normalizedFilter.DepartmentID != nil && *normalizedFilter.DepartmentID <= 0 {
		return EmployeeListResult{}, ErrInvalidInput
	}
	if strings.TrimSpace(normalizedFilter.Status) != "" {
		normalizedFilter.Status = strings.TrimSpace(normalizedFilter.Status)
	}

	rows, total, err := s.repo.ListEmployees(ctx, normalizedFilter)
	if err != nil {
		return EmployeeListResult{}, err
	}

	items := make([]EmployeeView, 0, len(rows))
	for _, row := range rows {
		items = append(items, ToEmployeeView(row))
	}

	return EmployeeListResult{
		Items:    items,
		Total:    total,
		Page:     normalizedFilter.Page,
		PageSize: normalizedFilter.PageSize,
	}, nil
}

func (s *Service) ListDepartments(ctx context.Context) ([]DepartmentOption, error) {
	rows, err := s.repo.ListDepartments(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]DepartmentOption, 0, len(rows))
	for _, row := range rows {
		items = append(items, DepartmentOption{
			ID:   row.ID,
			Name: row.Name,
		})
	}
	return items, nil
}

func (s *Service) normalizeInput(input UpsertEmployeeInput) (UpsertEmployeeInput, error) {
	normalized := input
	normalized.FirstName = strings.TrimSpace(input.FirstName)
	normalized.LastName = strings.TrimSpace(input.LastName)
	normalized.OtherName = strings.TrimSpace(input.OtherName)
	normalized.Gender = strings.TrimSpace(input.Gender)
	normalized.Phone = strings.TrimSpace(input.Phone)
	normalized.Email = strings.TrimSpace(input.Email)
	normalized.NationalID = strings.TrimSpace(input.NationalID)
	normalized.Address = strings.TrimSpace(input.Address)
	normalized.Position = strings.TrimSpace(input.Position)
	normalized.EmploymentStatus = strings.TrimSpace(input.EmploymentStatus)
	normalized.DOB = strings.TrimSpace(input.DOB)
	normalized.HireDate = strings.TrimSpace(input.HireDate)

	if normalized.FirstName == "" || normalized.LastName == "" {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}
	if normalized.Gender == "" || normalized.Phone == "" || normalized.Position == "" || normalized.EmploymentStatus == "" {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}
	if normalized.BaseSalary < 0 {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}
	if _, err := time.Parse("2006-01-02", normalized.DOB); err != nil {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}
	if _, err := time.Parse("2006-01-02", normalized.HireDate); err != nil {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}
	if normalized.Email != "" {
		if _, err := mail.ParseAddress(normalized.Email); err != nil {
			return UpsertEmployeeInput{}, ErrInvalidInput
		}
	}
	if normalized.DepartmentID != nil && *normalized.DepartmentID <= 0 {
		return UpsertEmployeeInput{}, ErrInvalidInput
	}

	return normalized, nil
}

func (s *Service) ensureDepartmentIntegrity(ctx context.Context, departmentID *int64) error {
	if departmentID == nil {
		return nil
	}
	exists, err := s.repo.DepartmentExists(ctx, *departmentID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrDepartmentNotFound
	}
	return nil
}
