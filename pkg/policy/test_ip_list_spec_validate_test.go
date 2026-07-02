//ff:func feature=policy type=engine control=sequence
//ff:what TestIPListSpec_Validate — tests validation of IPListSpec required Purpose field
package policy

import "testing"

func TestIPListSpec_Validate(t *testing.T) {
	if err := (&IPListSpec{Purpose: ""}).Validate(); err == nil {
		t.Fatal("expected error for empty Purpose")
	}
	if err := (&IPListSpec{Purpose: "blocklist"}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
