//ff:func feature=price type=engine control=sequence
//ff:what TestDiscountSpec_Validate — covers Rate range and Fixed non-negative branches of DiscountSpec.Validate
package price

import "testing"

func TestDiscountSpec_Validate(t *testing.T) {
	if err := (&DiscountSpec{Rate: -0.1}).Validate(); err == nil {
		t.Fatal("expected error for Rate below 0")
	}
	if err := (&DiscountSpec{Rate: 1.1}).Validate(); err == nil {
		t.Fatal("expected error for Rate above 1")
	}
	if err := (&DiscountSpec{Rate: 0.5, Fixed: -1}).Validate(); err == nil {
		t.Fatal("expected error for negative Fixed")
	}
	if err := (&DiscountSpec{Rate: 0.5, Fixed: 5}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
