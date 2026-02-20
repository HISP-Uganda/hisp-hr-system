package employees

import (
	"database/sql"
	"time"
)

type Employee struct {
	ID               int64          `db:"id" json:"id"`
	FirstName        string         `db:"first_name" json:"first_name"`
	LastName         string         `db:"last_name" json:"last_name"`
	OtherName        sql.NullString `db:"other_name" json:"-"`
	Gender           string         `db:"gender" json:"gender"`
	DOB              time.Time      `db:"dob" json:"-"`
	Phone            string         `db:"phone" json:"phone"`
	Email            sql.NullString `db:"email" json:"-"`
	NationalID       sql.NullString `db:"national_id" json:"-"`
	Address          sql.NullString `db:"address" json:"-"`
	DepartmentID     sql.NullInt64  `db:"department_id" json:"-"`
	DepartmentName   sql.NullString `db:"department_name" json:"-"`
	Position         string         `db:"position" json:"position"`
	EmploymentStatus string         `db:"employment_status" json:"employment_status"`
	HireDate         time.Time      `db:"hire_date" json:"-"`
	BaseSalary       float64        `db:"base_salary" json:"base_salary"`
	CreatedAt        time.Time      `db:"created_at" json:"-"`
	UpdatedAt        time.Time      `db:"updated_at" json:"-"`
}

type UpsertEmployeeInput struct {
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	OtherName        string  `json:"other_name"`
	Gender           string  `json:"gender"`
	DOB              string  `json:"dob"`
	Phone            string  `json:"phone"`
	Email            string  `json:"email"`
	NationalID       string  `json:"national_id"`
	Address          string  `json:"address"`
	DepartmentID     *int64  `json:"department_id"`
	Position         string  `json:"position"`
	EmploymentStatus string  `json:"employment_status"`
	HireDate         string  `json:"hire_date"`
	BaseSalary       float64 `json:"base_salary"`
}

type EmployeeListFilter struct {
	Search       string `json:"search"`
	DepartmentID *int64 `json:"department_id"`
	Status       string `json:"status"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
}

type EmployeeListResult struct {
	Items    []EmployeeView `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type EmployeeView struct {
	ID               int64   `json:"id"`
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	OtherName        string  `json:"other_name"`
	Gender           string  `json:"gender"`
	DOB              string  `json:"dob"`
	Phone            string  `json:"phone"`
	Email            string  `json:"email"`
	NationalID       string  `json:"national_id"`
	Address          string  `json:"address"`
	DepartmentID     *int64  `json:"department_id"`
	DepartmentName   string  `json:"department_name"`
	Position         string  `json:"position"`
	EmploymentStatus string  `json:"employment_status"`
	HireDate         string  `json:"hire_date"`
	BaseSalary       float64 `json:"base_salary"`
	CreatedAt        string  `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
}

type Department struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

type DepartmentOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func ToEmployeeView(row Employee) EmployeeView {
	var departmentID *int64
	if row.DepartmentID.Valid {
		value := row.DepartmentID.Int64
		departmentID = &value
	}

	return EmployeeView{
		ID:               row.ID,
		FirstName:        row.FirstName,
		LastName:         row.LastName,
		OtherName:        nullString(row.OtherName),
		Gender:           row.Gender,
		DOB:              row.DOB.Format("2006-01-02"),
		Phone:            row.Phone,
		Email:            nullString(row.Email),
		NationalID:       nullString(row.NationalID),
		Address:          nullString(row.Address),
		DepartmentID:     departmentID,
		DepartmentName:   nullString(row.DepartmentName),
		Position:         row.Position,
		EmploymentStatus: row.EmploymentStatus,
		HireDate:         row.HireDate.Format("2006-01-02"),
		BaseSalary:       row.BaseSalary,
		CreatedAt:        row.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt:        row.UpdatedAt.UTC().Format(time.RFC3339),
	}
}

func nullString(value sql.NullString) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
