package payroll

import "testing"

func TestCalculateAmounts(t *testing.T) {
	gross, net := CalculateAmounts(1000, 250, 80, 40)
	if gross != 1250 {
		t.Fatalf("expected gross 1250, got %v", gross)
	}
	if net != 1130 {
		t.Fatalf("expected net 1130, got %v", net)
	}
}
