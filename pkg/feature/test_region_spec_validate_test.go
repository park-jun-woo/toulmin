//ff:func feature=feature type=engine control=sequence
//ff:what TestRegionSpec_Validate — tests validation of RegionSpec region presence
package feature

import "testing"

func TestRegionSpec_Validate(t *testing.T) {
	if err := (&RegionSpec{Region: ""}).Validate(); err == nil {
		t.Fatal("expected error for empty region")
	}
	if err := (&RegionSpec{Region: "KR"}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
