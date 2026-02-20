ALTER TABLE payroll_entries
    DROP COLUMN IF EXISTS gross_pay;

ALTER TABLE payroll_entries
    RENAME COLUMN net_pay TO net_salary;

ALTER TABLE payroll_entries
    RENAME COLUMN tax_total TO tax;

ALTER TABLE payroll_entries
    RENAME COLUMN deductions_total TO deductions;

ALTER TABLE payroll_entries
    RENAME COLUMN allowances_total TO allowances;

ALTER TABLE payroll_entries
    RENAME COLUMN base_salary TO basic_salary;

DROP INDEX IF EXISTS idx_payroll_batches_month;

ALTER TABLE payroll_batches
    ADD COLUMN month_number INTEGER,
    ADD COLUMN year_number INTEGER,
    ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();

UPDATE payroll_batches
SET year_number = split_part(month, '-', 1)::INTEGER,
    month_number = split_part(month, '-', 2)::INTEGER;

ALTER TABLE payroll_batches
    ALTER COLUMN month_number SET NOT NULL,
    ALTER COLUMN year_number SET NOT NULL;

ALTER TABLE payroll_batches
    DROP CONSTRAINT IF EXISTS uq_payroll_batches_month,
    DROP CONSTRAINT IF EXISTS chk_payroll_batches_month_format;

ALTER TABLE payroll_batches
    DROP COLUMN IF EXISTS created_by,
    DROP COLUMN IF EXISTS approved_by,
    DROP COLUMN IF EXISTS approved_at,
    DROP COLUMN IF EXISTS locked_at,
    DROP COLUMN IF EXISTS month;

ALTER TABLE payroll_batches
    RENAME COLUMN month_number TO month;

ALTER TABLE payroll_batches
    RENAME COLUMN year_number TO year;

ALTER TABLE payroll_batches
    ADD CONSTRAINT uq_payroll_batches_month_year UNIQUE (month, year),
    ADD CONSTRAINT chk_payroll_batches_month_valid CHECK (month BETWEEN 1 AND 12);

CREATE INDEX idx_payroll_batches_year_month ON payroll_batches(year, month);
