# HISP HR System
## Authoritative Requirements Specification
## Phase A – Online-First Architecture
## Stack: Wails v2 (Go 1.22+) + React (TypeScript) + PostgreSQL (>= 13)

---

# 1. Vision

Build a secure, LAN-friendly, desktop HR Management System using Wails that:

- Runs as a desktop application (Windows/Mac/Linux)
- Connects to a centralized PostgreSQL database (online-first)
- Uses JWT authentication + role-based access control (RBAC)
- Manages Employees, Departments, Leave, and Payroll
- Is modular, scalable, and enterprise-ready
- Can later evolve to offline-sync (Phase B)

---

# 2. Target Architecture (Phase A)

## 2.1 High-Level

Frontend:
- React + TypeScript
- Material UI (MUI v5+)
    - Use dataGrid with advanced features (sticky headers, enable/disable/hide columns)
    - Support a Light and Dark theme, and a theme switcher
    - Use a good Snackbar/notification system
- TanStack Router
- TanStack Query

Backend:
- Go 1.22+
- Wails v2
- SQLX for database access
- golang-migrate for migrations
- JWT Authentication (access + refresh tokens)
- Structured logging (zerolog or logrus)
- Add tests whenever sensible

Database:
- PostgreSQL >= 13

Deployment:
- Single centralized PostgreSQL database (hosted on a server reachable on LAN/VPN)
- Desktop app distributed as a binary installer/package

---

# 3. Core Modules (Phase A)

## 3.1 Authentication (JWT)
- Login screen (small centered window)
- Username + Password
- On login:
  - Issue JWT access token (short-lived)
  - Issue refresh token (longer-lived)
- Refresh endpoint/function rotates or reissues access token
- Logout invalidates refresh token (server-side)
- RBAC enforced server-side on every protected operation

Roles:
- Admin (full access)
- HR Officer (employees, departments, leave)
- Finance Officer (payroll)
- Viewer (read-only where applicable)

---

## 3.2 Employee Management

CRUD fields:
- First Name
- Last Name
- Other Name (optional)
- Gender
- Date of Birth
- Phone
- Email (optional)
- National ID (optional)
- Address (optional)
- Position
- Department
- Employment Status
- Date of Hire
- Base Salary Amount

Search:
- By name
- By department
- By status

---

## 3.3 Department Management
- Create department
- Edit department
- Delete department
- Assign employees to departments
- Deletion rule (MVP):
  - A department with employees MUST NOT be deleted (enforce in backend)

---

## 3.4 Leave Management

### Leave Types
- Annual Leave
- Sick Leave
- Maternity Leave
- Unpaid Leave
- Custom types configurable

### Features
- Apply for leave
- Approve / Reject leave (HR/Admin)
- Leave status: Pending, Approved, Rejected
- Leave balance tracking
- Leave history per employee
- Automatic leave balance deduction upon approval

### Leave Rules (MVP defaults)
- Entitlement tracked per calendar year
- Leave cannot exceed available balance (based on approved leave days)

---

## 3.5 Payroll Management

### Payroll Features
- Monthly payroll processing using batches
- Base salary per employee (from employee record)
- Allowances (numeric)
- Deductions (numeric)
- Tax (numeric or percentage-based later; MVP uses numeric)
- Gross + Net calculations server-side

### Payroll Workflow
- Create batch (Draft)
- Generate payroll entries for all active employees (transactional)
- Edit entry values while Draft
- Approve batch (Approved)
- Lock batch (Locked; no further edits)
- Export payroll batch to CSV

RBAC:
- Finance Officer and Admin can manage payroll
- Others: no access (or read-only if explicitly enabled later)

---

## 3.6 User Management (Admin Only)
- Create system users
- Assign role
- Activate/deactivate accounts
- Reset passwords

---

## 3.7 Audit Logging (Phase A)
Audit events MUST be recorded for:
- user.login.success / user.login.fail
- token.refresh
- leave.request.create
- leave.request.approve / leave.request.reject
- payroll.batch.create / approve / lock
- payroll.entry.update
- user.create / user.deactivate / user.reset_password

Audit logs must include:
- actor_user_id
- action
- entity_type (optional)
- entity_id (optional)
- metadata JSON (optional)
- created_at timestamp

---

# 4. Security Requirements

- Password hashing with bcrypt
- JWT access tokens + refresh tokens
- Refresh tokens MUST be stored server-side and revocable
- RBAC checks enforced server-side for every protected operation
- Input validation server-side
- SQL injection protection (use SQLX with parameterized queries)
- No hardcoded credentials
- Sensitive config via environment variables or config file (not committed)

---

# 5. Database Requirements

PostgreSQL >= 13

## 5.1 Tables

### users
- id (PK)
- username (unique, not null)
- password_hash (not null)
- role (not null)
- is_active (not null, default true)
- created_at
- updated_at

### refresh_tokens
- id (PK)
- user_id (FK -> users.id, not null)
- token_hash (unique, not null)  # store hash, not raw token
- expires_at (not null)
- revoked_at (nullable)
- created_at
Indexes:
- user_id
- expires_at
- token_hash unique

### departments
- id (PK)
- name (unique, not null)
- description (nullable)
- created_at
- updated_at

### employees
- id (PK)
- first_name (not null)
- last_name (not null)
- other_name (nullable)
- gender (not null)
- dob (not null)
- phone (not null)
- email (nullable)
- national_id (nullable)
- address (nullable)
- department_id (FK -> departments.id, nullable initially allowed)
- position (not null)
- employment_status (not null)
- hire_date (not null)
- base_salary (not null, numeric)
- created_at
- updated_at
Indexes:
- (last_name, first_name)
- department_id
- employment_status

### leave_types
- id (PK)
- name (unique, not null)
- annual_entitlement_days (not null)
- is_active (not null, default true)
- created_at
- updated_at

### leave_requests
- id (PK)
- employee_id (FK -> employees.id, not null)
- leave_type_id (FK -> leave_types.id, not null)
- start_date (not null)
- end_date (not null)
- days_requested (not null)
- status (not null) # Pending|Approved|Rejected
- approved_by (FK -> users.id, nullable)
- approved_at (nullable)
- created_at
- updated_at
Indexes:
- employee_id
- status
- (start_date, end_date)

### payroll_batches
- id (PK)
- month (not null)
- year (not null)
- status (not null) # Draft|Approved|Locked
- created_at
- updated_at
Constraints:
- unique (month, year)
Indexes:
- (year, month)

### payroll_entries
- id (PK)
- batch_id (FK -> payroll_batches.id, not null)
- employee_id (FK -> employees.id, not null)
- basic_salary (not null, numeric)
- allowances (not null, numeric, default 0)
- deductions (not null, numeric, default 0)
- tax (not null, numeric, default 0)
- net_salary (not null, numeric)
- created_at
- updated_at
Constraints:
- unique (batch_id, employee_id)
Indexes:
- batch_id
- employee_id

### audit_logs
- id (PK)
- actor_user_id (FK -> users.id, nullable for system actions)
- action (not null)
- entity_type (nullable)
- entity_id (nullable)
- metadata (jsonb, nullable)
- created_at
Indexes:
- actor_user_id
- action
- created_at

---

# 6. Migrations Requirements

- Use golang-migrate
- Maintain up/down migration pairs
- Migrations live under: backend/migrations/
- App startup must run migrations (or provide a CLI command) for dev environments
- Production may run migrations via admin command/tooling (documented)

---

# 7. UI Requirements

## 7.1 Login Screen
- Small window (not much larger than form)
- Centered
- Clean professional design (MUI)
- No main shell visible before login

## 7.2 Main Shell
After authentication:
- Sidebar navigation
- Top bar with user info (username + role)
- Logout button
- Protected routes

Pages:
- Dashboard
- Employees
- Departments
- Leave
- Payroll
- Users (Admin only)
- (Optional) Audit Logs (Admin only)

---

# 8. Project Structure (Target)

```
hisp-hr-system/
  frontend/                 # React + TS (Wails template)
  backend/
    cmd/                    # optional CLI entrypoints
    internal/
      auth/
      users/
      employees/
      departments/
      leave/
      payroll/
      audit/
      db/
      middleware/
      config/
    migrations/
  docs/
    requirements.md
    codex-prompts.md
    status.md
```

---

# 9. Non-Functional Requirements

- Must compile on Windows, Mac, Linux
- Clean architecture and modular packages
- Structured logging across backend
- Database operations use SQLX with parameterized queries
- Payroll operations MUST be transactional
- Consistent error handling with user-friendly messages
- Basic automated tests for core services (auth, leave approval, payroll locking)

---

# 10. Out of Scope (Phase A)

- Offline mode / sync engine
- Recruitment
- Performance appraisal
- Mobile app
- External integrations

---

# 11. Future Phase B

- Offline-first sync engine
- Local SQLite cache
- Background sync service
- Conflict resolution strategy

---

# END OF AUTHORITATIVE SPEC
