# HISP HR System
## Development Status Tracker
## Phase A – Online-First (JWT + SQLX + golang-migrate)

Last Updated: 2026-02-20 06:06:29 UTC

---

# 1. Context Recovery Summary

Current implementation includes:
- Wails v2 desktop app bootstrapped and runnable
- Backend foundation (config + SQLX + migrations)
- Backend JWT auth service with refresh token storage/rotation
- Frontend login flow + auth state bootstrap + protected shell
- Main shell layout with role-aware sidebar visibility and route guards

Not yet implemented:
- Business modules (Employees/Departments/Leave/Payroll/User Management/Audit)
- Structured logging integration
- Seed/admin user creation path
- TanStack Query data layer usage

---

# 2. Implemented Architecture (Precise)

## Backend runtime wiring
- `app.go` initializes runtime at startup and stores DB/auth facade handles.
- `app_auth.go` exposes Wails bindings:
  - `Login(username, password)`
  - `Refresh(refreshToken)`
  - `Logout(refreshToken)`
  - `Me(accessToken)`
- `backend/bootstrap/foundation.go` loads config, runs migrations, opens SQLX DB, initializes auth facade.

## Config/DB/Migrations
- `backend/internal/config/config.go`
  - Environment-driven config (`APP_DB_URL`, `APP_JWT_SECRET`, token TTLs, DB pool tuning, migration flags/path).
- `backend/internal/db/db.go`
  - SQLX connection setup using `pgx` stdlib and pool settings.
- `backend/internal/db/migrate.go`
  - `golang-migrate` runner for `file://` migrations path.
- `backend/migrations/000001_init_schema.*.sql`
  - Full schema created for all required Phase A tables:
    - `users`, `refresh_tokens`, `departments`, `employees`, `leave_types`, `leave_requests`, `payroll_batches`, `payroll_entries`, `audit_logs`.

## Auth domain
- `backend/internal/auth/*`
  - SQLX auth repository (user lookup, refresh token insert/revoke/query-for-rotation).
  - bcrypt hash/verify helpers.
  - JWT access token generation/parsing (claims include `user_id`, `username`, `role`, expiry).
  - Refresh token generation (random token + SHA-256 hash persisted).
  - Login/Refresh/Logout service workflows.
- `backend/internal/middleware/jwt.go`
  - Token validation + user load + active-user check + auth context population.
- `backend/internal/middleware/rbac.go`
  - Role enforcement helper (`RequireRoles`).

## Frontend auth and routing
- `frontend/src/auth/AuthContext.tsx`
  - Auth bootstrap on app load (`Me`, fallback to `Refresh`, otherwise clear session).
  - Login/logout actions via Wails bindings.
  - Session tokens persisted in localStorage.
- `frontend/src/routes/router.tsx`
  - Protected route tree with authenticated shell parent.
  - Routes:
    - `/dashboard`
    - `/employees`
    - `/departments`
    - `/leave`
    - `/payroll`
    - `/users`
  - Unauthorized route access redirects to `/dashboard`.
  - Unauthenticated access redirects to `/login`.
- `frontend/src/components/AppShell.tsx`
  - Responsive sidebar + top bar + logout.
  - Role-based menu visibility.

---

# 3. Completed Milestones

✔ Phase 1 Foundation complete
- Config loader, SQLX DB setup, migrations, startup migration runner.

✔ Phase 2 Authentication backend complete
- JWT + refresh tokens, bcrypt, repository/service/middleware, Wails auth bindings.

✔ Phase 3 Login UI + auth state complete
- Centered login form (MUI), auth bootstrap, protected rendering.

✔ Phase 4 Main shell complete
- Authenticated shell, sidebar/top bar, route guards, role-aware navigation.

---

# 4. Verified Build State

- Backend: `go test ./...` passes (current packages compile; limited tests exist).
- Frontend: `cd frontend && npm run build` passes.

---

# 5. Gaps and Risks (Current)

## Functional gaps
- No seeded/default user provisioning currently; login requires manually inserting a user with bcrypt hash.
- Module backends (employees/departments/leave/payroll/users/audit actions) not yet implemented.
- Route pages are placeholders; no module CRUD UI yet.

## Technical debt
- Structured logging (zerolog/logrus) not integrated.
- TanStack Query not yet integrated for data fetching/caching.
- Auth token storage is currently localStorage-based (not in-memory only).
- Backend RBAC middleware exists but has not yet been integrated across module handlers/services (because modules are not yet built).

## Testing gaps
- No integration tests for login/refresh/logout against a real DB.
- No middleware/RBAC behavior tests.

---

# 6. Pending Work by Milestone

1. Phase 5 Employees Module
- Backend CRUD + search/pagination + frontend table/dialogs.

2. Phase 6 Departments Module
- CRUD + safe delete rule (prevent delete with assigned employees).

3. Phase 7 Leave Module
- Leave types + leave requests + approval/balance logic.

4. Phase 8 Payroll Module
- Batch lifecycle + transactional entry generation + CSV export.

5. Phase 9 User Management
- Admin user CRUD, activate/deactivate, reset password.

6. Phase 10 Audit Logging
- Event logger integration for required critical actions.

7. Phase 11 Hardening
- RBAC/validation/error handling audit, tests, cross-platform build checks.

---

# 7. Next Immediate Action

➡ Execute Codex Prompt #5 (Employees Module)
Goal: Deliver first business module end-to-end (backend CRUD/search + frontend table/dialogs).

---

# END
