//ff:func feature=state type=engine control=sequence
//ff:what TestExpirySpec_SpecName — verifies ExpirySpec.SpecName returns fixed name
package state

import "testing"

func TestExpirySpec_SpecName(t *testing.T) {
	if got := (&ExpirySpec{}).SpecName(); got != "ExpirySpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "ExpirySpec")
	}
}
