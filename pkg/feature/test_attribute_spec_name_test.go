//ff:func feature=feature type=engine control=sequence
//ff:what TestAttributeSpec_SpecName — tests SpecName returns the fixed spec name
package feature

import "testing"

func TestAttributeSpec_SpecName(t *testing.T) {
	spec := &AttributeSpec{Key: "beta", Value: true}
	if got := spec.SpecName(); got != "AttributeSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "AttributeSpec")
	}
}
