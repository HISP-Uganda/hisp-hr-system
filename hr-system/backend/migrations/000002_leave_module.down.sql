ALTER TABLE leave_requests
    DROP CONSTRAINT IF EXISTS chk_leave_requests_working_days_positive;

ALTER TABLE leave_requests
    DROP COLUMN IF EXISTS comment,
    DROP COLUMN IF EXISTS cancelled_at,
    DROP COLUMN IF EXISTS cancelled_by,
    DROP COLUMN IF EXISTS rejected_at,
    DROP COLUMN IF EXISTS rejected_by,
    DROP COLUMN IF EXISTS requested_by,
    DROP COLUMN IF EXISTS working_days;

DROP INDEX IF EXISTS idx_leave_locked_dates_lock_date;
DROP TABLE IF EXISTS leave_locked_dates;

DROP INDEX IF EXISTS idx_leave_entitlements_employee_year;
DROP TABLE IF EXISTS leave_entitlements;

ALTER TABLE employees
    DROP COLUMN IF EXISTS user_id;

ALTER TABLE leave_types
    DROP COLUMN IF EXISTS counts_toward_entitlement,
    DROP COLUMN IF EXISTS requires_approval,
    DROP COLUMN IF EXISTS requires_attachment,
    DROP COLUMN IF EXISTS is_paid;
