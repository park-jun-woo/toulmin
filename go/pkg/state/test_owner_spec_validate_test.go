//ff:func feature=state type=engine control=sequence
//ff:what TestOwnerSpec_Validate — tests OwnerSpec.Validate always returns nil
package state

import "testing"

func TestOwnerSpec_Validate(t *testing.T) {
	err := (&OwnerSpec{}).Validate()
	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}
