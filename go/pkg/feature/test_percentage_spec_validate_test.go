//ff:func feature=feature type=engine control=sequence
//ff:what TestPercentageSpec_Validate — tests validation of PercentageSpec percentage bounds
package feature

import "testing"

func TestPercentageSpec_Validate(t *testing.T) {
	if err := (&PercentageSpec{Percentage: -0.1}).Validate(); err == nil {
		t.Fatal("expected error for percentage below zero")
	}
	if err := (&PercentageSpec{Percentage: 1.1}).Validate(); err == nil {
		t.Fatal("expected error for percentage above one")
	}
	if err := (&PercentageSpec{Percentage: 0.5}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
