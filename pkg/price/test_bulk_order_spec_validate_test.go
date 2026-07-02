//ff:func feature=price type=engine control=sequence
//ff:what TestBulkOrderSpec_Validate — tests validation of BulkOrderSpec positive MinQuantity
package price

import "testing"

func TestBulkOrderSpec_Validate(t *testing.T) {
	if err := (&BulkOrderSpec{MinQuantity: 0}).Validate(); err == nil {
		t.Fatal("expected error for zero MinQuantity")
	}
	if err := (&BulkOrderSpec{MinQuantity: -1}).Validate(); err == nil {
		t.Fatal("expected error for negative MinQuantity")
	}
	if err := (&BulkOrderSpec{MinQuantity: 5}).Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
