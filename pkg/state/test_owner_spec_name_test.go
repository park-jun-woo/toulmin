//ff:func feature=state type=engine control=sequence
//ff:what TestOwnerSpec_SpecName — tests OwnerSpec.SpecName returns fixed name
package state

import "testing"

func TestOwnerSpec_SpecName(t *testing.T) {
	name := (&OwnerSpec{}).SpecName()
	if name != "OwnerSpec" {
		t.Errorf("expected %q, got %q", "OwnerSpec", name)
	}
}
