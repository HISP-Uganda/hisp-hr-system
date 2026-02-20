package payroll

import "time"

const (
	StatusDraft    = "Draft"
	StatusApproved = "Approved"
	StatusLocked   = "Locked"
)

type Actor struct {
	UserID int64
	Role   string
}

type Batch struct {
	ID         int64      `db:"id" json:"id"`
	Month      string     `db:"month" json:"month"`
	Status     string     `db:"status" json:"status"`
	CreatedBy  int64      `db:"created_by" json:"created_by"`
	CreatedAt  time.Time  `db:"created_at" json:"created_at"`
	ApprovedBy *int64     `db:"approved_by" json:"approved_by,omitempty"`
	ApprovedAt *time.Time `db:"approved_at" json:"approved_at,omitempty"`
	LockedAt   *time.Time `db:"locked_at" json:"locked_at,omitempty"`
}

type Entry struct {
	ID              int64     `db:"id" json:"id"`
	BatchID         int64     `db:"batch_id" json:"batch_id"`
	EmployeeID      int64     `db:"employee_id" json:"employee_id"`
	EmployeeName    string    `db:"employee_name" json:"employee_name"`
	BaseSalary      float64   `db:"base_salary" json:"base_salary"`
	AllowancesTotal float64   `db:"allowances_total" json:"allowances_total"`
	DeductionsTotal float64   `db:"deductions_total" json:"deductions_total"`
	TaxTotal        float64   `db:"tax_total" json:"tax_total"`
	GrossPay        float64   `db:"gross_pay" json:"gross_pay"`
	NetPay          float64   `db:"net_pay" json:"net_pay"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type BatchFilter struct {
	Month  string `json:"month"`
	Status string `json:"status"`
}

type BatchDetail struct {
	Batch   Batch   `json:"batch"`
	Entries []Entry `json:"entries"`
}

type CreateBatchInput struct {
	Month string `json:"month"`
}

type UpdateEntryAmountsInput struct {
	AllowancesTotal float64 `json:"allowances_total"`
	DeductionsTotal float64 `json:"deductions_total"`
	TaxTotal        float64 `json:"tax_total"`
}

type BatchListResult struct {
	Items []Batch `json:"items"`
}
