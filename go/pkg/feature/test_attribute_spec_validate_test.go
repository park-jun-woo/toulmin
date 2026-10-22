//ff:func feature=feature type=engine control=sequence
//ff:what TestAttributeSpec_Validate — tests validation of AttributeSpec key presence
package feature

import "testing"

func TestAttributeSpec_Validate(t *testing.T) {
	if err := (&AttributeSpec{Key: "", Value: true}).Validate(); err == nil {
		t.Fatal("expected error for empty key")
	}
	if err := (&AttributeSpec{Key: "beta", Value: true}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
