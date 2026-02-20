ALTER TABLE payroll_batches
    RENAME COLUMN month TO month_number;

ALTER TABLE payroll_batches
    RENAME COLUMN year TO year_number;

ALTER TABLE payroll_batches
    DROP CONSTRAINT IF EXISTS uq_payroll_batches_month_year;

ALTER TABLE payroll_batches
    DROP CONSTRAINT IF EXISTS chk_payroll_batches_month_valid;

DROP INDEX IF EXISTS idx_payroll_batches_year_month;

ALTER TABLE payroll_batches
    ADD COLUMN month TEXT,
    ADD COLUMN created_by BIGINT REFERENCES users(id),
    ADD COLUMN approved_by BIGINT REFERENCES users(id),
    ADD COLUMN approved_at TIMESTAMPTZ,
    ADD COLUMN locked_at TIMESTAMPTZ;

UPDATE payroll_batches
SET month = to_char(make_date(year_number, month_number, 1), 'YYYY-MM')
WHERE month IS NULL;

ALTER TABLE payroll_batches
    ALTER COLUMN month SET NOT NULL;

ALTER TABLE payroll_batches
    ADD CONSTRAINT uq_payroll_batches_month UNIQUE (month);

ALTER TABLE payroll_batches
    ADD CONSTRAINT chk_payroll_batches_month_format CHECK (month ~ '^[0-9]{4}-(0[1-9]|1[0-2])$');

ALTER TABLE payroll_batches
    DROP COLUMN year_number,
    DROP COLUMN month_number,
    DROP COLUMN updated_at;

CREATE INDEX idx_payroll_batches_month ON payroll_batches(month);

ALTER TABLE payroll_entries
    RENAME COLUMN basic_salary TO base_salary;

ALTER TABLE payroll_entries
    RENAME COLUMN allowances TO allowances_total;

ALTER TABLE payroll_entries
    RENAME COLUMN deductions TO deductions_total;

ALTER TABLE payroll_entries
    RENAME COLUMN tax TO tax_total;

ALTER TABLE payroll_entries
    RENAME COLUMN net_salary TO net_pay;

ALTER TABLE payroll_entries
    ADD COLUMN gross_pay NUMERIC(14,2) NOT NULL DEFAULT 0;

UPDATE payroll_entries
SET gross_pay = base_salary + allowances_total,
    net_pay = (base_salary + allowances_total) - deductions_total - tax_total;

ALTER TABLE payroll_entries
    ALTER COLUMN allowances_total SET NOT NULL,
    ALTER COLUMN deductions_total SET NOT NULL,
    ALTER COLUMN tax_total SET NOT NULL,
    ALTER COLUMN net_pay SET NOT NULL,
    ALTER COLUMN gross_pay SET NOT NULL;

ALTER TABLE payroll_entries
    ALTER COLUMN gross_pay DROP DEFAULT;
