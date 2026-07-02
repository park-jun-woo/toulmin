//ff:func feature=policy type=engine control=sequence
//ff:what TestIPListSpec_SpecName — verifies IPListSpec.SpecName returns fixed name
package policy

import "testing"

func TestIPListSpec_SpecName(t *testing.T) {
	if got := (&IPListSpec{}).SpecName(); got != "IPListSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "IPListSpec")
	}
}
