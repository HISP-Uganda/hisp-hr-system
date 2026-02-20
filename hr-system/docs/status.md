# HISP HR System
## Development Status Tracker
## Phase A – Online-First (JWT + SQLX + golang-migrate)

Last Updated: 2026-02-20 09:19:33 UTC

---

# 1. Context Recovery Summary

Implemented and working:
- Wails v2 desktop runtime with Go backend + React frontend.
- PostgreSQL connectivity via SQLX and startup migrations via golang-migrate.
- JWT auth with refresh-token rotation and server-side revocation.
- Role-aware app shell and protected routing.
- Employees module (backend CRUD/search + frontend list/dialog workflows).
- Leave module (backend business rules + frontend planner/apply/requests/balance flows).

Not implemented yet:
- Departments module (Phase 6)
- Payroll module (Phase 8)
- User management module (Phase 9)
- Audit logging module (Phase 10)
- Structured logging integration
- TanStack Query integration
- Seed/admin bootstrap workflow

---

# 2. Implemented Architecture (Precise)

## Backend runtime and layering
- Runtime wiring in `app.go` + `backend/bootstrap/foundation.go`.
- Active facades:
  - Auth: `backend/bootstrap/auth.go`
  - Employees: `backend/bootstrap/employees.go`
  - Leave: `backend/bootstrap/leave.go`
- Bound Wails app APIs:
  - Auth bindings: `app_auth.go`
  - Employee bindings: `app_employees.go`
  - Leave bindings: `app_leave.go`
- Clean architecture followed per module:
  - Wails app method -> facade/service -> repository -> SQLX/Postgres

## Database and migrations
- Base schema: `backend/migrations/000001_init_schema.*.sql`
- Leave extension migration: `backend/migrations/000002_leave_module.*.sql`
  - Adds leave policy fields to `leave_types`
  - Adds `leave_entitlements`
  - Adds `leave_locked_dates`
  - Extends `leave_requests` with `working_days` and status metadata fields
  - Adds `employees.user_id` for self-service leave ownership resolution

## Auth module (complete)
- JWT access/refresh flow with hashed refresh tokens in DB.
- Middleware + role checks:
  - `backend/internal/middleware/jwt.go`
  - `backend/internal/middleware/rbac.go`

## Employees module (complete)
- Backend:
  - `backend/internal/employees/repository.go`
  - `backend/internal/employees/service.go`
- Features:
  - Create/Update/Delete/Get/List
  - Name/department/status filtering with pagination
  - Server-side validation
  - Department integrity checks on create/update
- Frontend:
  - `frontend/src/modules/employees/EmployeesPage.tsx`
  - Table + search + filters + create/edit/delete dialogs

## Leave module (complete, implemented ahead of Department milestone by user request)
- Backend domain:
  - `backend/internal/leave/repository.go`
  - `backend/internal/leave/service.go`
  - `backend/internal/leave/rules.go`
  - `backend/internal/leave/errors.go`
- Core rules implemented:
  - Working-days calculation server-side (weekends excluded)
  - Reject invalid ranges / zero working days
  - Reject locked-date collisions
  - Reject overlaps with approved leave for same employee
  - Balance formula enforcement:
    - `available = entitlement_total - reserved - (pending + approved)`
  - Status transitions:
    - Pending -> Approved (HR/Admin)
    - Pending -> Rejected (HR/Admin)
    - Pending -> Cancelled (self or HR/Admin)
    - Approved -> Cancelled (HR/Admin)
  - Master-only edit/delete paths (role strings `Master` / `Master Admin`)
- Leave UI:
  - `frontend/src/modules/leave/LeavePage.tsx`
  - Planner tab (locked dates)
  - Apply form with working-days preview
  - Requests/history with approve/reject/cancel controls
  - Balance report view
- API mapping notes:
  - `docs/notes/leave.md` records mapping of requested REST semantics to Wails bindings.

---

# 3. Milestone Status

Completed:
1. Phase 1 Foundation
2. Phase 2 Authentication
3. Phase 3 Login UI + auth state
4. Phase 4 Main shell
5. Phase 5 Employees module
6. Phase 7 Leave module (out of sequence by explicit request)

Not started:
1. Phase 6 Departments module
2. Phase 8 Payroll module
3. Phase 9 User management
4. Phase 10 Audit logging
5. Phase 11 hardening

---

# 4. Verified Build/Test State

Most recent known verification:
- Backend: `GOCACHE=$(pwd)/.gocache go test ./...` passes.
- Frontend: `cd frontend && npm run build` passes.
- Leave unit tests present:
  - `backend/internal/leave/rules_test.go`
  - `backend/internal/leave/service_test.go`

---

# 5. Risks and Gaps

Functional gaps:
- Department CRUD and safe-delete enforcement still missing.
- Attendance module is not implemented; leave conversion path contains TODO for full attendance state mutation.

Testing gaps:
- No DB-backed repository integration tests for leave/employees.
- No end-to-end UI workflow tests.

Technical debt:
- Structured logging not integrated.
- TanStack Query not used for frontend data layer yet.
- Auth tokens still persisted via localStorage.

---

# 6. Next Work Item (from status order)

Next NOT STARTED item:
- Phase 6 Departments module
  - Backend CRUD
  - Enforce delete protection when employees are assigned
  - Frontend departments list + dialogs

---

# END
