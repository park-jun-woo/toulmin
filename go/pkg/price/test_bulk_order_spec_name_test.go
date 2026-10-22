//ff:func feature=price type=engine control=sequence
//ff:what TestBulkOrderSpec_SpecName — verifies BulkOrderSpec.SpecName returns fixed name
package price

import "testing"

func TestBulkOrderSpec_SpecName(t *testing.T) {
	if got := (&BulkOrderSpec{}).SpecName(); got != "BulkOrderSpec" {
		t.Fatalf("SpecName() = %q, want %q", got, "BulkOrderSpec")
	}
}
