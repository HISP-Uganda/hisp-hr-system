ALTER TABLE leave_types
    ADD COLUMN IF NOT EXISTS is_paid BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS requires_attachment BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS counts_toward_entitlement BOOLEAN NOT NULL DEFAULT TRUE;

ALTER TABLE employees
    ADD COLUMN IF NOT EXISTS user_id BIGINT UNIQUE REFERENCES users(id);

CREATE TABLE IF NOT EXISTS leave_entitlements (
    id BIGSERIAL PRIMARY KEY,
    employee_id BIGINT NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    leave_type_id BIGINT NOT NULL REFERENCES leave_types(id) ON DELETE CASCADE,
    year INTEGER NOT NULL,
    total_days INTEGER NOT NULL,
    reserved_days INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_leave_entitlements_employee_type_year UNIQUE (employee_id, leave_type_id, year),
    CONSTRAINT chk_leave_entitlements_year CHECK (year >= 2000),
    CONSTRAINT chk_leave_entitlements_total_nonnegative CHECK (total_days >= 0),
    CONSTRAINT chk_leave_entitlements_reserved_nonnegative CHECK (reserved_days >= 0)
);
CREATE INDEX IF NOT EXISTS idx_leave_entitlements_employee_year ON leave_entitlements(employee_id, year);

CREATE TABLE IF NOT EXISTS leave_locked_dates (
    id BIGSERIAL PRIMARY KEY,
    lock_date DATE NOT NULL UNIQUE,
    reason TEXT,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_leave_locked_dates_lock_date ON leave_locked_dates(lock_date);

ALTER TABLE leave_requests
    ADD COLUMN IF NOT EXISTS working_days INTEGER,
    ADD COLUMN IF NOT EXISTS requested_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS rejected_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS rejected_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS cancelled_by BIGINT REFERENCES users(id),
    ADD COLUMN IF NOT EXISTS cancelled_at TIMESTAMPTZ,
    ADD COLUMN IF NOT EXISTS comment TEXT;

UPDATE leave_requests
SET working_days = days_requested
WHERE working_days IS NULL;

ALTER TABLE leave_requests
    ALTER COLUMN working_days SET NOT NULL;

ALTER TABLE leave_requests
    ADD CONSTRAINT chk_leave_requests_working_days_positive CHECK (working_days > 0);
