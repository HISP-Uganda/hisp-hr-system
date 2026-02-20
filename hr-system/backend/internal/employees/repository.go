package employees

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateEmployee(ctx context.Context, input UpsertEmployeeInput) (Employee, error) {
	const query = `
		INSERT INTO employees (
			first_name,
			last_name,
			other_name,
			gender,
			dob,
			phone,
			email,
			national_id,
			address,
			department_id,
			position,
			employment_status,
			hire_date,
			base_salary
		)
		VALUES (
			:first_name,
			:last_name,
			:other_name,
			:gender,
			:dob,
			:phone,
			:email,
			:national_id,
			:address,
			:department_id,
			:position,
			:employment_status,
			:hire_date,
			:base_salary
		)
		RETURNING
			id,
			first_name,
			last_name,
			other_name,
			gender,
			dob,
			phone,
			email,
			national_id,
			address,
			department_id,
			NULL::TEXT AS department_name,
			position,
			employment_status,
			hire_date,
			base_salary,
			created_at,
			updated_at
	`

	row := mapToNamedParams(input)
	rows, err := r.db.NamedQueryContext(ctx, query, row)
	if err != nil {
		return Employee{}, fmt.Errorf("create employee: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return Employee{}, ErrEmployeeNotFound
	}

	var employee Employee
	if err := rows.StructScan(&employee); err != nil {
		return Employee{}, fmt.Errorf("scan created employee: %w", err)
	}
	return employee, nil
}

func (r *Repository) UpdateEmployee(ctx context.Context, employeeID int64, input UpsertEmployeeInput) (Employee, error) {
	const query = `
		UPDATE employees
		SET
			first_name = :first_name,
			last_name = :last_name,
			other_name = :other_name,
			gender = :gender,
			dob = :dob,
			phone = :phone,
			email = :email,
			national_id = :national_id,
			address = :address,
			department_id = :department_id,
			position = :position,
			employment_status = :employment_status,
			hire_date = :hire_date,
			base_salary = :base_salary,
			updated_at = NOW()
		WHERE id = :id
		RETURNING
			id,
			first_name,
			last_name,
			other_name,
			gender,
			dob,
			phone,
			email,
			national_id,
			address,
			department_id,
			NULL::TEXT AS department_name,
			position,
			employment_status,
			hire_date,
			base_salary,
			created_at,
			updated_at
	`

	row := mapToNamedParams(input)
	row["id"] = employeeID
	rows, err := r.db.NamedQueryContext(ctx, query, row)
	if err != nil {
		return Employee{}, fmt.Errorf("update employee: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return Employee{}, ErrEmployeeNotFound
	}

	var employee Employee
	if err := rows.StructScan(&employee); err != nil {
		return Employee{}, fmt.Errorf("scan updated employee: %w", err)
	}
	return employee, nil
}

func (r *Repository) DeleteEmployee(ctx context.Context, employeeID int64) error {
	const query = `DELETE FROM employees WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, employeeID)
	if err != nil {
		return fmt.Errorf("delete employee: %w", err)
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("read delete result: %w", err)
	}
	if affected == 0 {
		return ErrEmployeeNotFound
	}
	return nil
}

func (r *Repository) GetEmployee(ctx context.Context, employeeID int64) (Employee, error) {
	const query = `
		SELECT
			e.id,
			e.first_name,
			e.last_name,
			e.other_name,
			e.gender,
			e.dob,
			e.phone,
			e.email,
			e.national_id,
			e.address,
			e.department_id,
			d.name AS department_name,
			e.position,
			e.employment_status,
			e.hire_date,
			e.base_salary,
			e.created_at,
			e.updated_at
		FROM employees e
		LEFT JOIN departments d ON d.id = e.department_id
		WHERE e.id = $1
	`
	var row Employee
	if err := r.db.GetContext(ctx, &row, query, employeeID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Employee{}, ErrEmployeeNotFound
		}
		return Employee{}, fmt.Errorf("get employee: %w", err)
	}
	return row, nil
}

func (r *Repository) ListEmployees(ctx context.Context, filter EmployeeListFilter) ([]Employee, int, error) {
	whereClause, args := buildEmployeeFilters(filter)

	countQuery := `
		SELECT COUNT(1)
		FROM employees e
		LEFT JOIN departments d ON d.id = e.department_id
	` + whereClause

	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args...); err != nil {
		return nil, 0, fmt.Errorf("count employees: %w", err)
	}

	listArgs := append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)
	listQuery := `
		SELECT
			e.id,
			e.first_name,
			e.last_name,
			e.other_name,
			e.gender,
			e.dob,
			e.phone,
			e.email,
			e.national_id,
			e.address,
			e.department_id,
			d.name AS department_name,
			e.position,
			e.employment_status,
			e.hire_date,
			e.base_salary,
			e.created_at,
			e.updated_at
		FROM employees e
		LEFT JOIN departments d ON d.id = e.department_id
	` + whereClause + `
		ORDER BY e.last_name ASC, e.first_name ASC, e.id ASC
		LIMIT $` + fmt.Sprintf("%d", len(listArgs)-1) + ` OFFSET $` + fmt.Sprintf("%d", len(listArgs))

	rows := make([]Employee, 0)
	if err := r.db.SelectContext(ctx, &rows, listQuery, listArgs...); err != nil {
		return nil, 0, fmt.Errorf("list employees: %w", err)
	}

	return rows, total, nil
}

func (r *Repository) DepartmentExists(ctx context.Context, departmentID int64) (bool, error) {
	const query = `SELECT id FROM departments WHERE id = $1`
	var department Department
	if err := r.db.GetContext(ctx, &department, query, departmentID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check department exists: %w", err)
	}
	return true, nil
}

func (r *Repository) ListDepartments(ctx context.Context) ([]Department, error) {
	const query = `
		SELECT id, name
		FROM departments
		ORDER BY name ASC
	`
	rows := make([]Department, 0)
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, fmt.Errorf("list departments: %w", err)
	}
	return rows, nil
}

func buildEmployeeFilters(filter EmployeeListFilter) (string, []any) {
	conditions := make([]string, 0, 3)
	args := make([]any, 0, 3)
	position := 1

	search := strings.TrimSpace(filter.Search)
	if search != "" {
		pattern := "%" + search + "%"
		conditions = append(conditions,
			"(e.first_name ILIKE $"+fmt.Sprintf("%d", position)+" OR e.last_name ILIKE $"+fmt.Sprintf("%d", position)+" OR COALESCE(e.other_name, '') ILIKE $"+fmt.Sprintf("%d", position)+")",
		)
		args = append(args, pattern)
		position++
	}

	if filter.DepartmentID != nil {
		conditions = append(conditions, "e.department_id = $"+fmt.Sprintf("%d", position))
		args = append(args, *filter.DepartmentID)
		position++
	}

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		conditions = append(conditions, "e.employment_status = $"+fmt.Sprintf("%d", position))
		args = append(args, status)
		position++
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func mapToNamedParams(input UpsertEmployeeInput) map[string]any {
	return map[string]any{
		"first_name":        input.FirstName,
		"last_name":         input.LastName,
		"other_name":        nullableString(input.OtherName),
		"gender":            input.Gender,
		"dob":               input.DOB,
		"phone":             input.Phone,
		"email":             nullableString(input.Email),
		"national_id":       nullableString(input.NationalID),
		"address":           nullableString(input.Address),
		"department_id":     nullableInt64(input.DepartmentID),
		"position":          input.Position,
		"employment_status": input.EmploymentStatus,
		"hire_date":         input.HireDate,
		"base_salary":       input.BaseSalary,
	}
}

func nullableString(value string) any {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return trimmed
}

func nullableInt64(value *int64) any {
	if value == nil || *value <= 0 {
		return nil
	}
	return *value
}
