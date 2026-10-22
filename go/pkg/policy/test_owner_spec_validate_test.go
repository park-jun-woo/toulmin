//ff:func feature=policy type=engine control=sequence
//ff:what TestOwnerSpec_Validate — verifies OwnerSpec.Validate always returns nil
package policy

import "testing"

func TestOwnerSpec_Validate(t *testing.T) {
	if err := (&OwnerSpec{}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
