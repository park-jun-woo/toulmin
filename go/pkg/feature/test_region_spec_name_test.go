//ff:func feature=feature type=engine control=sequence
//ff:what TestRegionSpec_SpecName — tests SpecName returns the fixed spec name
package feature

import "testing"

func TestRegionSpec_SpecName(t *testing.T) {
	spec := &RegionSpec{Region: "KR"}
	if got := spec.SpecName(); got != "RegionSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "RegionSpec")
	}
}
