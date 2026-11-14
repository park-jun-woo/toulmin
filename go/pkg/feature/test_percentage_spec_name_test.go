//ff:func feature=feature type=engine control=sequence
//ff:what TestPercentageSpec_SpecName — tests SpecName returns the fixed spec name
package feature

import "testing"

func TestPercentageSpec_SpecName(t *testing.T) {
	spec := &PercentageSpec{Percentage: 0.5}
	if got := spec.SpecName(); got != "PercentageSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "PercentageSpec")
	}
}
