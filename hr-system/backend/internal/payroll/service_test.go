package payroll

import (
	"context"
	"testing"
	"time"
)

type fakeStore struct {
	batches map[int64]Batch
	entries map[int64]Entry

	generateCalls int
	approveCalls  int
	lockCalls     int
}

func (f *fakeStore) ListBatches(_ context.Context, _ BatchFilter) ([]Batch, error) {
	items := make([]Batch, 0, len(f.batches))
	for _, batch := range f.batches {
		items = append(items, batch)
	}
	return items, nil
}

func (f *fakeStore) GetBatch(_ context.Context, batchID int64) (Batch, error) {
	item, ok := f.batches[batchID]
	if !ok {
		return Batch{}, ErrBatchNotFound
	}
	return item, nil
}

func (f *fakeStore) GetBatchEntries(_ context.Context, batchID int64) ([]Entry, error) {
	items := make([]Entry, 0)
	for _, entry := range f.entries {
		if entry.BatchID == batchID {
			items = append(items, entry)
		}
	}
	return items, nil
}

func (f *fakeStore) GetEntry(_ context.Context, entryID int64) (Entry, error) {
	item, ok := f.entries[entryID]
	if !ok {
		return Entry{}, ErrEntryNotFound
	}
	return item, nil
}

func (f *fakeStore) CreateBatch(_ context.Context, month string, createdBy int64) (Batch, error) {
	id := int64(len(f.batches) + 1)
	batch := Batch{ID: id, Month: month, Status: StatusDraft, CreatedBy: createdBy, CreatedAt: time.Now().UTC()}
	f.batches[id] = batch
	return batch, nil
}

func (f *fakeStore) GenerateEntriesForBatch(_ context.Context, _ int64) error {
	f.generateCalls++
	return nil
}

func (f *fakeStore) UpdateEntryAmounts(_ context.Context, entryID int64, allowancesTotal, deductionsTotal, taxTotal, grossPay, netPay float64) (Entry, error) {
	entry, ok := f.entries[entryID]
	if !ok {
		return Entry{}, ErrEntryNotFound
	}
	entry.AllowancesTotal = allowancesTotal
	entry.DeductionsTotal = deductionsTotal
	entry.TaxTotal = taxTotal
	entry.GrossPay = grossPay
	entry.NetPay = netPay
	f.entries[entryID] = entry
	return entry, nil
}

func (f *fakeStore) ApproveBatch(_ context.Context, batchID int64, approvedBy int64, approvedAt time.Time) (Batch, error) {
	batch := f.batches[batchID]
	batch.Status = StatusApproved
	batch.ApprovedBy = &approvedBy
	batch.ApprovedAt = &approvedAt
	f.batches[batchID] = batch
	f.approveCalls++
	return batch, nil
}

func (f *fakeStore) LockBatch(_ context.Context, batchID int64, lockedAt time.Time) (Batch, error) {
	batch := f.batches[batchID]
	batch.Status = StatusLocked
	batch.LockedAt = &lockedAt
	f.batches[batchID] = batch
	f.lockCalls++
	return batch, nil
}

func newTestService() *Service {
	store := &fakeStore{
		batches: map[int64]Batch{
			1: {ID: 1, Month: "2026-02", Status: StatusDraft, CreatedBy: 1, CreatedAt: time.Now().UTC()},
			2: {ID: 2, Month: "2026-01", Status: StatusApproved, CreatedBy: 1, CreatedAt: time.Now().UTC()},
		},
		entries: map[int64]Entry{
			10: {
				ID:              10,
				BatchID:         1,
				EmployeeID:      21,
				EmployeeName:    "Doe, Jane",
				BaseSalary:      1000,
				AllowancesTotal: 0,
				DeductionsTotal: 0,
				TaxTotal:        0,
				GrossPay:        1000,
				NetPay:          1000,
			},
			11: {
				ID:              11,
				BatchID:         2,
				EmployeeID:      22,
				EmployeeName:    "Doe, John",
				BaseSalary:      1200,
				AllowancesTotal: 10,
				DeductionsTotal: 2,
				TaxTotal:        1,
				GrossPay:        1210,
				NetPay:          1207,
			},
		},
	}
	svc, _ := NewService(store)
	return svc
}

func TestStatusTransitions(t *testing.T) {
	svc := newTestService()

	if _, err := svc.ApproveBatch(context.Background(), Actor{UserID: 9, Role: "Finance Officer"}, 1); err != nil {
		t.Fatalf("approve draft batch: %v", err)
	}
	if _, err := svc.LockBatch(context.Background(), Actor{UserID: 9, Role: "Finance Officer"}, 1); err != nil {
		t.Fatalf("lock approved batch: %v", err)
	}
	if _, err := svc.ApproveBatch(context.Background(), Actor{UserID: 9, Role: "Finance Officer"}, 2); err != ErrInvalidStatusTransition {
		t.Fatalf("expected invalid status transition, got %v", err)
	}
}

func TestUpdateEntryOnlyDraftBatch(t *testing.T) {
	svc := newTestService()

	_, err := svc.UpdateEntryAmounts(context.Background(), Actor{UserID: 9, Role: "Finance Officer"}, 10, UpdateEntryAmountsInput{
		AllowancesTotal: 100,
		DeductionsTotal: 30,
		TaxTotal:        20,
	})
	if err != nil {
		t.Fatalf("update draft entry failed: %v", err)
	}

	_, err = svc.UpdateEntryAmounts(context.Background(), Actor{UserID: 9, Role: "Finance Officer"}, 11, UpdateEntryAmountsInput{
		AllowancesTotal: 100,
		DeductionsTotal: 30,
		TaxTotal:        20,
	})
	if err != ErrBatchImmutable {
		t.Fatalf("expected immutable error for approved batch entry, got %v", err)
	}
}
