//ff:func feature=policy type=engine control=sequence
//ff:what TestOwnerSpec_SpecName — verifies OwnerSpec.SpecName returns fixed name
package policy

import "testing"

func TestOwnerSpec_SpecName(t *testing.T) {
	if got := (&OwnerSpec{}).SpecName(); got != "OwnerSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "OwnerSpec")
	}
}
