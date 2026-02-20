# Leave Module Notes

## Architecture Decision
- The authoritative system architecture is Wails backend + React frontend.
- Requested HTTP endpoint contracts were mapped to Wails-bound backend methods with the same functional behavior:
  - `GET /me/leave/balance?year=YYYY` -> `MeLeaveBalance(accessToken, year)`
  - `GET /admin/leave/balance/:employeeId?year=YYYY` -> `AdminLeaveBalance(accessToken, employeeID, year)`
  - `GET /admin/leave/requests` -> `ListLeaveRequests(accessToken, filter)`
  - `POST /leave/apply` -> `ApplyLeave(accessToken, input)`
  - `POST /leave/:id/approve` -> `ApproveLeave(accessToken, id, input)`
  - `POST /leave/:id/reject` -> `RejectLeave(accessToken, id, input)`
  - `POST /leave/:id/cancel` -> `CancelLeave(accessToken, id, input)`

## Schema
- Migration: `backend/migrations/000002_leave_module.up.sql`
- Added/updated structures:
  - Extended `leave_types` with:
    - `is_paid`, `requires_attachment`, `requires_approval`, `counts_toward_entitlement`
  - Added `leave_entitlements`
  - Added `leave_locked_dates`
  - Extended `leave_requests` with:
    - `working_days`, `requested_by`, `rejected_by`, `rejected_at`, `cancelled_by`, `cancelled_at`, `comment`
  - Added nullable `employees.user_id` for self-service resolution

## Validation Rules Implemented
- Working days computed server-side from date range excluding weekends.
- Invalid date ranges rejected (`end < start`, no working days).
- Locked-date collision rejected when any working date is locked.
- Overlap with approved leave for same employee rejected.
- Available balance check:
  - `available = entitlement_total - reserved - (pending + approved)`
- Status transitions:
  - `Pending -> Approved` (Admin/HR)
  - `Pending -> Rejected` (Admin/HR)
  - `Pending -> Cancelled` (employee self or Admin/HR)
  - `Approved -> Cancelled` (Admin/HR)

## Master Operations
- Master-only (`Master` / `Master Admin`) operations:
  - Edit any leave request (`MasterUpdateLeave`)
  - Delete any leave request (`MasterDeleteLeave`)

## Attendance Integration
- Implemented `ConvertAbsenceToLeave` flow as an application-level bridge that creates one-day approved-path leave request validation flow.
- Attendance state mutation is marked TODO in service code until attendance module exists.
