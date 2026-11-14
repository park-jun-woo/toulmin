//ff:func feature=policy type=engine control=sequence
//ff:what TestHeaderSpec_Validate — tests validation of HeaderSpec required Header field
package policy

import "testing"

func TestHeaderSpec_Validate(t *testing.T) {
	if err := (&HeaderSpec{Header: ""}).Validate(); err == nil {
		t.Fatal("expected error for empty Header")
	}
	if err := (&HeaderSpec{Header: "X-Test"}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
