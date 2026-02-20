package payroll

import (
	"context"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func TestGenerateEntriesTransactionalRollback(t *testing.T) {
	dsn := os.Getenv("PAYROLL_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("PAYROLL_TEST_DATABASE_URL is not set")
	}

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer db.Close()

	ctx := context.Background()
	setup := []string{
		`DROP TABLE IF EXISTS payroll_entries`,
		`DROP TABLE IF EXISTS payroll_batches`,
		`DROP TABLE IF EXISTS employees`,
		`CREATE TABLE employees (id BIGINT PRIMARY KEY, employment_status TEXT NOT NULL, base_salary NUMERIC(14,2) NOT NULL)`,
		`CREATE TABLE payroll_batches (id BIGSERIAL PRIMARY KEY, month TEXT NOT NULL, status TEXT NOT NULL, created_by BIGINT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), approved_by BIGINT, approved_at TIMESTAMPTZ, locked_at TIMESTAMPTZ)`,
		`CREATE TABLE payroll_entries (id BIGSERIAL PRIMARY KEY, batch_id BIGINT NOT NULL, employee_id BIGINT NOT NULL, base_salary NUMERIC(14,2) NOT NULL, allowances_total NUMERIC(14,2) NOT NULL, deductions_total NUMERIC(14,2) NOT NULL, tax_total NUMERIC(14,2) NOT NULL, gross_pay NUMERIC(14,2) NOT NULL, net_pay NUMERIC(14,2) NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW())`,
		`INSERT INTO payroll_batches (id, month, status, created_by) VALUES (1, '2026-02', 'Draft', 1)`,
		`INSERT INTO employees (id, employment_status, base_salary) VALUES (1, 'Active', 1000), (2, 'Active', 1200)`,
		`CREATE OR REPLACE FUNCTION fail_second_payroll_insert() RETURNS trigger AS $$ BEGIN IF NEW.employee_id = 2 THEN RAISE EXCEPTION 'boom'; END IF; RETURN NEW; END; $$ LANGUAGE plpgsql`,
		`CREATE TRIGGER payroll_entries_fail_second BEFORE INSERT ON payroll_entries FOR EACH ROW EXECUTE FUNCTION fail_second_payroll_insert()`,
	}
	for _, stmt := range setup {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			t.Fatalf("setup failed on %q: %v", stmt, err)
		}
	}
	defer func() {
		_, _ = db.ExecContext(ctx, `DROP TRIGGER IF EXISTS payroll_entries_fail_second ON payroll_entries`)
		_, _ = db.ExecContext(ctx, `DROP FUNCTION IF EXISTS fail_second_payroll_insert`)
		_, _ = db.ExecContext(ctx, `DROP TABLE IF EXISTS payroll_entries`)
		_, _ = db.ExecContext(ctx, `DROP TABLE IF EXISTS payroll_batches`)
		_, _ = db.ExecContext(ctx, `DROP TABLE IF EXISTS employees`)
	}()

	repo := NewRepository(db)
	err = repo.GenerateEntriesForBatch(ctx, 1)
	if err == nil {
		t.Fatalf("expected generation failure")
	}

	var count int
	if err := db.GetContext(ctx, &count, `SELECT COUNT(1) FROM payroll_entries WHERE batch_id = 1`); err != nil {
		t.Fatalf("count entries: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected rollback to leave zero entries, got %d", count)
	}
}
