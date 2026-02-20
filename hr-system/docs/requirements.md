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

# 3.4 Leave Management

## 3.4.1 Overview

The system shall provide a complete leave management module that supports leave applications, approvals, entitlement tracking, calendar planning, attendance integration, and reporting.

Leave management must support a centralized, role-based workflow with automatic leave balance calculations and validation.

---

## 3.4.2 Leave Types

### Default Leave Types

The system shall include the following default leave types:

- Annual Leave
- Sick Leave
- Maternity Leave
- Study Leave
- Unpaid Leave

### Configurable Leave Types

Administrators shall be able to:

- Create new leave types
- Edit existing leave types
- Deactivate leave types

Each leave type shall support the following properties:

- Name
- Paid / Unpaid (boolean)
- Requires Attachment (boolean)
- Requires Approval (boolean)
- Counts Toward Annual Entitlement (boolean)
- Active (boolean)

---

## 3.4.3 Leave Entitlements

Each employee shall have defined leave entitlements per leave year.

### Leave Year

- Default leave year: Calendar year (January–December)
- Future support for configurable leave year

### Entitlement Fields Per Employee

- Total Leave Days (Annual entitlement)
- Reserved Leave Days (Non-usable buffer)
- Year

### Leave Balance Calculation

Available Leave Days shall be calculated as:

Available = Total Entitlement - Reserved Leave - (Approved Leave Days + Pending Leave Days)

Rules:

- Leave request cannot exceed available balance.
- Leave balance recalculates automatically when:
  - Leave is approved
  - Leave is rejected
  - Leave is cancelled
  - Leave record is deleted
  - Entitlement is modified

---

## 3.4.4 Leave Request Workflow

### Employee Capabilities

Employees shall be able to:

- Apply for leave
- Select:
  - Leave type
  - Start date
  - End date
- View automatic working day calculation
- View real-time leave balance validation
- Edit leave request if status is Pending
- Cancel leave request if status is Pending
- View personal leave history

### Automatic Calculations

When selecting Start Date and End Date, the system shall:

1. Validate End Date ≥ Start Date
2. Calculate working days only (exclude weekends)
3. Exclude locked/restricted dates
4. Prevent submission if:
   - Insufficient leave balance
   - Date range contains no working days
   - Restricted dates are included
   - Overlapping approved leave exists

---

## 3.4.5 Leave Status Lifecycle

Each leave request shall have one of the following statuses:

- Pending
- Approved
- Rejected
- Cancelled

### Status Transitions

| From     | To        | Performed By |
|----------|----------|--------------|
| Pending  | Approved | HR/Admin     |
| Pending  | Rejected | HR/Admin     |
| Pending  | Cancelled | Employee    |
| Approved | Cancelled | HR/Admin    |

Approved leave shall immediately reduce available leave balance.

Rejected or Cancelled leave shall restore leave balance.

---

## 3.4.6 Approval & Permissions

### HR / Admin

HR/Admin users shall be able to:

- View all leave requests
- Approve or Reject leave requests
- Create leave requests on behalf of employees
- Override leave balance (with justification)
- Lock or unlock specific calendar dates
- View department-level leave reports

### Master Admin

Master Admin shall additionally be able to:

- Edit any leave record (including approved)
- Delete leave records
- Modify entitlement values

### Staff Users

Staff users shall:

- Apply for leave
- View personal leave history
- Edit or cancel pending leave only
- View personal leave balance

---

## 3.4.7 Leave Planner (Annual Calendar View)

The system shall provide a year-based leave planner calendar.

Each date in the calendar shall have one of the following states:

- Available
- Locked (Admin controlled)
- Scheduled (Employee planning)

### Locked Dates

Locked dates:

- Cannot be selected for leave
- Are configurable in system settings
- May represent:
  - Public holidays
  - Organizational blackout periods
  - Critical operational dates

---

## 3.4.8 Attendance Integration

The system shall integrate leave with attendance records.

Capabilities:

- Convert an Absence record to Approved Leave (1 day deduction)
- Automatically update attendance status to “Leave”
- Prevent conversion if insufficient leave balance

---

## 3.4.9 Reporting

The system shall generate the following reports:

- Leave requests report (by period)
- Leave balances report
- Leave usage by department
- Leave history per employee
- Export to CSV
- Printable PDF output

Reports shall support filtering by:

- Department
- Date range
- Leave type
- Employee

---

## 3.4.10 Data Model Requirements

### LeaveRequest

Fields:

- id
- employee_id
- leave_type_id
- start_date
- end_date
- working_days
- status
- reason (optional)
- approved_by (nullable)
- approved_at (nullable)
- created_at
- updated_at

### LeaveType

Fields:

- id
- name
- paid (boolean)
- counts_toward_entitlement (boolean)
- requires_attachment (boolean)
- requires_approval (boolean)
- active (boolean)

### LeaveEntitlement

Fields:

- employee_id
- year
- total_days
- reserved_days

---

## 3.4.11 Future Enhancements (Non-MVP)

The architecture shall allow future support for:

- Leave carry-forward policy
- Monthly accrual-based leave
- Multi-level approval workflows
- Email notifications
- Attachment uploads (e.g., medical certificate)
- Team leave overlap warnings
- Departmental leave heatmap
- Leave quotas per department

---

## 3.5 Payroll Management

### 3.5.1 Overview

The system shall provide a structured, batch-based payroll engine that supports:

- Monthly payroll processing
- Transactional generation of payroll entries
- Controlled approval workflow
- Locking of finalized payroll
- Server-side calculation of gross and net pay
- Auditability and immutability after approval

Payroll must be processed in discrete monthly batches.

---

## 3.5.2 Payroll Batches

Payroll shall be processed per month using a Payroll Batch entity.

### PayrollBatch Fields

- id
- month (YYYY-MM format)
- status (Draft | Approved | Locked)
- created_by
- created_at
- approved_by (nullable)
- approved_at (nullable)
- locked_at (nullable)

### Batch Status Lifecycle

| From     | To        | Performed By |
|----------|----------|--------------|
| Draft    | Approved | Finance/Admin |
| Approved | Locked   | Finance/Admin |

Rules:

- Draft: Entries may be edited.
- Approved: Entries are finalized but visible.
- Locked: No further modifications allowed.
- Locked batches are immutable.

Only one batch per month shall be allowed.

---

## 3.5.3 Payroll Entries

When a batch is created, payroll entries shall be generated transactionally for all active employees.

### PayrollEntry Fields

- id
- batch_id
- employee_id
- base_salary
- allowances_total
- deductions_total
- tax_total
- gross_pay
- net_pay
- created_at
- updated_at

Entries are tied to a specific batch.

---

## 3.5.4 Payroll Calculation Rules

All payroll calculations must be performed server-side.

### Calculation Components

1. Base Salary  
   Retrieved from employee record.

2. Allowances (numeric total for MVP)
  - Stored per payroll entry
  - Editable while batch is Draft

3. Deductions (numeric total for MVP)
  - Stored per payroll entry
  - Editable while batch is Draft

4. Tax (numeric for MVP)
  - Manual numeric value
  - Future: percentage-based or bracket-based engine

### Gross Pay

Gross Pay = Base Salary + Allowances

### Net Pay

Net Pay = Gross Pay - Deductions - Tax


All computed values must be persisted in payroll_entries.

---

## 3.5.5 Payroll Generation (Transactional)

When generating payroll entries for a Draft batch:

- The operation must be transactional.
- All active employees must receive an entry.
- If any employee entry fails, the entire generation must roll back.
- Regeneration is allowed only while batch is Draft.

---

## 3.5.6 Editing Rules

While batch status is Draft:

- Allowances may be edited.
- Deductions may be edited.
- Tax value may be edited.
- Recalculation must occur server-side after edits.

After batch is Approved:

- No financial fields may be edited.
- Only viewing and exporting allowed.

After batch is Locked:

- No changes permitted.
- Entries are immutable.

---

## 3.5.7 Exporting Payroll

The system shall support exporting a payroll batch to CSV.

Export must include:

- Employee Name
- Base Salary
- Allowances
- Deductions
- Tax
- Gross Pay
- Net Pay

Export shall only be allowed for Approved or Locked batches.

---

## 3.5.8 RBAC (Role-Based Access Control)

### Finance Officer
- Create payroll batch
- Generate payroll entries
- Edit Draft batch
- Approve batch
- Lock batch
- Export payroll

### Admin
- Same permissions as Finance Officer

### Master Admin
- All Finance/Admin capabilities
- May unlock batch (optional if later enabled)
- May delete Draft batch

### Other Users
- No access to payroll module
- Future: optional read-only access if explicitly granted

---

## 3.5.9 Audit Requirements

The system shall:

- Record created_by and approved_by users
- Prevent modification of Approved/Locked batches
- Maintain payroll history per month
- Prevent duplicate batches per month

---

## 3.5.10 Future Enhancements (Non-MVP)

Architecture must allow future support for:

- Per-allowance breakdown (multiple allowance types)
- Per-deduction breakdown
- Automated tax engine
- Statutory contributions (NSSF, PAYE engine)
- Payslip generation (PDF)
- Bank transfer file export
- Project-based salary allocation
- Multi-currency payroll

---

## 3.6 User Management (Admin Only)

### Overview

This module allows administrators to manage system users.

Only users with role `admin` may access this functionality (both backend APIs and frontend UI).

---

## Backend Requirements

### Access Control
- All endpoints require:
    - Valid JWT
    - Role = `admin`

---

### Data Model

Minimum fields for `users` table:

- `id` (int, primary key)
- `username` (unique, indexed, required)
- `password_hash` (required)
- `role` (string, required)
- `is_active` (boolean, default true)
- `created_at` (timestamp)
- `updated_at` (timestamp)
- `last_login_at` (nullable timestamp)

---

### Endpoints

#### 1. Create User
- Admin only
- Input:
    - username
    - password
    - role
- Validation:
    - Username must be unique
    - Password must meet minimum length requirement
- Password must be hashed before storage
- Response must not include password

---

#### 2. Update User
- Admin only
- Allows updating:
    - username
    - role
- Must not update password here
- Returns updated user

---

#### 3. Get User
- Admin only
- Fetch single user by ID
- Must not return password

---

#### 4. List Users
- Admin only
- Supports:
    - Pagination (`page`, `pageSize`)
    - Optional search (`q`) by username
- Returns:
    - List of users
    - Total count

---

#### 5. Reset Password
- Admin only
- Input:
    - new password
- Must hash password before saving
- Return success response only

---

#### 6. Activate / Deactivate User
- Admin only
- Toggle `is_active`
- System must prevent an admin from deactivating their own account
- Return updated user

---

### Security Requirements

- Passwords must use a strong hashing algorithm (bcrypt or argon2).
- Passwords must never be returned in API responses.
- All user management actions must be logged in audit logs.

---

## Frontend Requirements (React + MUI)

### Access
- Users page visible only to `admin` role.

---

### Users Page

Features:
- Users table
- Pagination controls
- Search input (username)
- Columns:
    - Username
    - Role
    - Status (Active/Inactive)
    - Created At
    - Actions

---

### Create User Dialog
- Username field
- Password field
- Role selector
- Validation
- Submit button

---

### Edit User Dialog
- Edit username
- Change role
- Save button

---

### Reset Password Dialog
- New password
- Confirm password
- Validation
- Submit button

---

### Activate / Deactivate Control
- Toggle or action button
- Confirmation dialog before change

---

## Deliverables

- Admin can:
    - Create users
    - Assign/change roles
    - Reset passwords
    - Activate/deactivate users
    - Search and paginate users
    - Initial admin user seeded via environment configuration

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
