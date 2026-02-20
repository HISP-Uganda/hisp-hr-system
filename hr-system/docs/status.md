# HISP HR System
## Development Status Tracker
## Phase A – Online-First (JWT + SQLX + golang-migrate)

Last Updated: 2026-02-20 05:12:10 UTC

---

# 1. Project State

Authoritative specification:
- docs/requirements.md (JWT-based auth, SQLX, golang-migrate)

Technology stack confirmed:
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

---

# 2. Completed

✔ Phase A scope finalized
✔ Payroll included in Phase A
✔ Leave Management included in Phase A
✔ JWT authentication selected (access + refresh tokens)
✔ SQLX selected for DB access
✔ golang-migrate selected for migrations
✔ Logical database schema fully defined
✔ RBAC roles defined:
    - Admin
    - HR Officer
    - Finance Officer
    - Viewer
✔ Audit logging requirements defined
✔ Codex workflow prompts prepared
✔ Foundation phase executed (Codex Prompt #1)
✔ Authentication backend phase executed (Codex Prompt #2)
✔ Login UI + Auth State phase executed (Codex Prompt #3):
    - Added backend Wails auth bindings in `app_auth.go`:
        - `Login(username, password)`
        - `Refresh(refreshToken)`
        - `Logout(refreshToken)`
        - `Me(accessToken)`
    - Extended bootstrap runtime auth wiring:
        - `backend/bootstrap/foundation.go`
        - `backend/bootstrap/auth.go`
    - Added frontend auth state management:
        - `frontend/src/auth/AuthContext.tsx`
        - token storage + session bootstrap/refresh
    - Added small centered MUI login screen:
        - `frontend/src/components/LoginPage.tsx`
    - Added protected routing with TanStack Router:
        - `frontend/src/routes/router.tsx`
        - login route + authenticated route guard
    - Ensured shell is not rendered before auth:
        - `AppRouter` returns nothing while auth initializes
        - protected route redirects unauthenticated users to `/login`
    - Added authenticated shell placeholder with logout:
        - `frontend/src/components/ShellPage.tsx`
    - Updated Wails frontend bindings:
        - `frontend/wailsjs/go/main/App.js`
        - `frontend/wailsjs/go/main/App.d.ts`
    - Verification passed:
        - `go test ./...`
        - `cd frontend && npm run build`

---

# 3. In Progress

🚧 Phase 4: Main Shell

Planned implementation tasks:
- Sidebar navigation
- Top bar with username + role + logout
- Route-based menu visibility
- Hide unauthorized menu items
- Fully protected app routes

---

# 4. Pending Milestones

## Foundation
- Add structured logging (zerolog or logrus) across backend

## Main Shell
- Sidebar navigation
- Top bar with user info
- Route-based access visibility
- Hide unauthorized menu items

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

- Structured logging is not yet integrated in backend flows.
- TanStack Query is not yet wired into auth/module data fetching.

---

# 7. Known Risks

- Payroll rules may evolve (tax logic, allowances structure)
- Leave calculation edge cases (overlaps, cross-year leave)
- Token revocation logic must be carefully implemented
- RBAC consistency across all handlers
- Future offline extension will require architectural layering

---

# 8. Next Immediate Action

➡ Execute Codex Prompt #4 (Main Shell)
Goal: Build final authenticated shell layout with role-aware navigation and route/menu visibility.

---

# END
