# HISP HR System
## Development Status Tracker
## Phase A – Online-First (JWT + SQLX + golang-migrate)

Last Updated: 2026-02-19 12:24:14 UTC

---

# 1. Project State

Authoritative specification:
- docs/requirements.md (JWT-based auth, SQLX, golang-migrate)

Technology stack planned:
- Wails v2
- Go 1.22+
- React + TypeScript
- Material UI (MUI v5+)
- TanStack Router
- TanStack Query
- PostgreSQL >= 13
- SQLX
- golang-migrate
- JWT (access + refresh tokens)

Implementation has NOT started yet.

---

# 2. Completed

None.

---

# 3. In Progress

None.

---

# 4. Pending Milestones

## Foundation
- Initialize Wails project
- Implement config loader (DB, JWT secrets)
- Setup SQLX DB connection
- Integrate golang-migrate
- Create full migration set (all tables)
- Wire migration runner into startup
- Add structured logging

## Authentication (JWT)
- User model
- bcrypt password hashing
- JWT access token generation
- Refresh token rotation & storage (hashed)
- Token validation middleware
- RBAC middleware
- Login UI (small centered window)
- Auth state management (frontend)

## Main Shell
- Sidebar navigation
- Protected routes
- Role-based menu visibility
- "Me" endpoint/function

## Employees Module
- CRUD backend
- Search (name, department, status)
- Frontend table + dialogs

## Departments Module
- CRUD backend
- Enforce safe delete rule
- Department UI

## Leave Module
- leave_types CRUD
- leave_requests workflow
- Calendar-year entitlement logic
- Balance computation
- Approval workflow (HR/Admin)
- Leave UI

## Payroll Module
- payroll_batches workflow
- payroll_entries generation (transactional)
- Draft → Approved → Locked logic
- Net salary calculation server-side
- CSV export
- Payroll UI

## User Management
- Admin-only CRUD
- Activate/deactivate
- Reset password

## Audit Logging
- audit_logs table
- Event logging helpers
- Record critical actions

## Hardening & QA
- Input validation review
- RBAC enforcement review
- Error handling standardization
- Cross-platform build test
- Basic service-layer tests

---

# 5. Architectural Decisions

1. Online-first centralized PostgreSQL
2. JWT-based authentication with refresh tokens stored server-side (hashed)
3. SQLX for explicit SQL and performance control
4. golang-migrate for versioned schema management
5. Clean backend architecture:
   handlers → services → repositories → db
6. Payroll operations MUST be transactional
7. Leave balance tracked per calendar year
8. Department deletion prohibited if employees are assigned

---

# 6. Technical Debt

None (project not yet started).

---

# 7. Known Risks

- Payroll rules may evolve (tax logic, allowances structure)
- Leave calculation edge cases (overlaps, cross-year leave)
- Token revocation logic must be carefully implemented
- RBAC consistency across all handlers
- Future offline extension will require architectural layering

---

# 8. Next Immediate Action

➡ Execute Codex Phase 1 (Foundation)

---

# END
