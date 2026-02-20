package leave

import (
	"testing"
	"time"
)

func TestComputeWorkingDays(t *testing.T) {
	t.Run("same-day weekday", func(t *testing.T) {
		start := time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC) // Monday
		days, dates := ComputeWorkingDays(start, start)
		if days != 1 || len(dates) != 1 {
			t.Fatalf("expected 1 working day, got %d", days)
		}
	})

	t.Run("weekend excluded", func(t *testing.T) {
		start := time.Date(2026, 2, 20, 0, 0, 0, 0, time.UTC) // Friday
		end := time.Date(2026, 2, 23, 0, 0, 0, 0, time.UTC)   // Monday
		days, _ := ComputeWorkingDays(start, end)
		if days != 2 {
			t.Fatalf("expected 2 working days, got %d", days)
		}
	})

	t.Run("invalid range", func(t *testing.T) {
		start := time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC)
		end := time.Date(2026, 3, 9, 0, 0, 0, 0, time.UTC)
		days, _ := ComputeWorkingDays(start, end)
		if days != 0 {
			t.Fatalf("expected 0 working days, got %d", days)
		}
	})
}

func TestCalculateAvailableBalance(t *testing.T) {
	available := CalculateAvailableBalance(20, 2, 5, 4)
	if available != 9 {
		t.Fatalf("expected available 9, got %d", available)
	}

	available = CalculateAvailableBalance(10, 3, 4, 5)
	if available != 0 {
		t.Fatalf("expected available floor at 0, got %d", available)
	}
}
