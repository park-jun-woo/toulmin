//ff:func feature=state type=engine control=sequence
//ff:what TestExpirySpec_Validate — verifies ExpirySpec.Validate always returns nil
package state

import "testing"

func TestExpirySpec_Validate(t *testing.T) {
	if err := (&ExpirySpec{}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
