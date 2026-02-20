Before implementing anything:

1. Read docs/requirements.md
2. Read docs/status.md
3. Follow clean architecture:
   handlers → services → repositories → db
4. All DB access MUST use SQLX with parameterized queries.
5. All schema changes MUST use golang-migrate (up + down files).
6. JWT authentication with:
   - Short-lived access token
   - Refresh token stored hashed in DB
7. All protected operations MUST enforce RBAC server-side.
8. Payroll operations MUST be transactional.
9. Update docs/status.md after every milestone.
10. Do not introduce external frameworks unless specified.
