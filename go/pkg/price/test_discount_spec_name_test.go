//ff:func feature=price type=engine control=sequence
//ff:what TestDiscountSpec_SpecName — verifies DiscountSpec.SpecName returns fixed name
package price

import "testing"

func TestDiscountSpec_SpecName(t *testing.T) {
	if got := (&DiscountSpec{}).SpecName(); got != "DiscountSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "DiscountSpec")
	}
}
