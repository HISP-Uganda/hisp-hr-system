# User Management Module Notes

## Scope
Implemented against `docs/requirements.md` section `3.6` using Wails bindings (desktop API surface) instead of HTTP routes.

## Endpoint to Binding Mapping
All operations require `accessToken` and enforce server-side admin-only authorization.

- `POST /api/admin/users` -> `CreateUser(accessToken, input)`
  - `input`: `{ username, password, role }`
- `GET /api/admin/users?page=&pageSize=&q=` -> `ListUsers(accessToken, query)`
  - `query`: `{ page, page_size, q }`
- `GET /api/admin/users/:id` -> `GetUser(accessToken, userID)`
- `PUT /api/admin/users/:id` -> `UpdateUser(accessToken, userID, input)`
  - `input`: `{ username, role }`
- `POST /api/admin/users/:id/reset-password` -> `ResetUserPassword(accessToken, userID, input)`
  - `input`: `{ new_password }`
- `POST /api/admin/users/:id/status` -> `SetUserStatus(accessToken, userID, input)`
  - `input`: `{ is_active }`

## Response/Error Conventions
Wails responses use the existing app envelope for data-returning methods:
- `{ success, message, data }`

Error strings are normalized to status-coded semantics:
- `401 unauthorized`
- `403 forbidden`
- `404 not found`
- `409 conflict`
- `422 invalid input`

## Key Rules Enforced
- Valid JWT required on every operation.
- Role must be `admin` (case-insensitive check).
- Password hash is never returned in responses.
- Username must be unique (conflict on duplicate).
- Passwords are hashed with bcrypt on create and reset.
- Password reset and create require minimum length (8 chars).
- Admin cannot deactivate their own account.
- List supports pagination and username search.

## Data Model + Migration
Migration: `backend/migrations/000004_user_management.up.sql`
- Adds nullable `users.last_login_at`
- Adds `idx_users_username`

`users` fields used by module:
- `id`, `username`, `password_hash`, `role`, `is_active`, `created_at`, `updated_at`, `last_login_at`

## Audit Logging
User management actions are recorded in `audit_logs` with `entity_type='user'`:
- `user.create`
- `user.update`
- `user.reset_password`
- `user.activate`
- `user.deactivate`

## Seed Bootstrap
Initial admin seeding is idempotent via env configuration:
- `APP_INITIAL_ADMIN_USERNAME`
- `APP_INITIAL_ADMIN_PASSWORD`
- `APP_INITIAL_ADMIN_ROLE` (default `admin`)

Behavior:
- Insert-on-conflict by username (`ON CONFLICT DO NOTHING`).
