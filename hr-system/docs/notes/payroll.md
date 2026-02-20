# Payroll Module Notes

## Scope
Implemented against `docs/requirements.md` section `3.5` with Wails bindings (desktop API surface) instead of HTTP routes.

## Backend Binding Surface
All methods require `accessToken` and enforce server-side RBAC (`Admin` or `Finance Officer`).

- `ListPayrollBatches(accessToken, filter)`
  - `filter.month` (`YYYY-MM`, optional)
  - `filter.status` (`Draft|Approved|Locked`, optional)
- `GetPayrollBatch(accessToken, batchID)`
  - Returns batch and entries
- `CreatePayrollBatch(accessToken, { month })`
  - Month format: `YYYY-MM`
  - Fails on duplicate month
- `GeneratePayrollEntries(accessToken, batchID)`
  - Transactional regenerate while Draft (delete + recreate)
  - Populates active employees only
- `UpdatePayrollEntryAmounts(accessToken, entryID, { allowances_total, deductions_total, tax_total })`
  - Allowed only when parent batch is Draft
  - Recomputes and persists gross/net server-side
- `ApprovePayrollBatch(accessToken, batchID)`
  - Allowed only from Draft
  - Sets `approved_by`, `approved_at`, status `Approved`
- `LockPayrollBatch(accessToken, batchID)`
  - Allowed only from Approved
  - Sets `locked_at`, status `Locked`
- `ExportPayrollBatchCSV(accessToken, batchID)`
  - Allowed only when batch is Approved or Locked
  - CSV columns:
    - Employee Name
    - Base Salary
    - Allowances
    - Deductions
    - Tax
    - Gross Pay
    - Net Pay

## Data Model Alignment
Migration: `backend/migrations/000003_payroll_module.up.sql`

- `payroll_batches`
  - `id`
  - `month` (`YYYY-MM`)
  - `status` (`Draft|Approved|Locked`)
  - `created_by`, `created_at`
  - `approved_by`, `approved_at` (nullable)
  - `locked_at` (nullable)
  - unique month enforced
- `payroll_entries`
  - `id`
  - `batch_id`
  - `employee_id`
  - `base_salary`
  - `allowances_total`
  - `deductions_total`
  - `tax_total`
  - `gross_pay` (persisted, computed server-side)
  - `net_pay` (persisted, computed server-side)
  - `created_at`, `updated_at`

## Calculation Rules
Server-side and persisted:

- `gross_pay = base_salary + allowances_total`
- `net_pay = gross_pay - deductions_total - tax_total`

## Status and Immutability Rules
- Draft:
  - generation/regeneration allowed
  - entry amount edits allowed
  - approval allowed
- Approved:
  - no edits/regeneration
  - lock allowed
  - CSV export allowed
- Locked:
  - immutable
  - CSV export allowed

## Frontend Screens
- `PayrollBatchesPage`
  - list month/status
  - create batch
  - open details
- `PayrollBatchDetailPage`
  - entries table
  - inline amount edits only in Draft
  - Generate/Approve/Lock/Export controls shown by status
  - status + timestamps shown

## Tests
- Unit: calculation correctness
  - `backend/internal/payroll/calculation_test.go`
- Unit: status transition and edit guards
  - `backend/internal/payroll/service_test.go`
- Integration-style repository test: transactional rollback on generation failure
  - `backend/internal/payroll/repository_integration_test.go`
  - uses `PAYROLL_TEST_DATABASE_URL`
