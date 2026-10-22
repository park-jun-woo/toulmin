//ff:func feature=tangl type=parser control=sequence
//ff:what TestBumpCounter — the Run/Check counter selected by slotKind is the one advanced
package en

import "testing"

// TestBumpCounter checks that slotKind 1 (Run) advances runN and leaves checkN
// untouched, and any other slotKind (Check, 2) advances checkN and leaves runN.
func TestBumpCounter(t *testing.T) {
	runN, checkN := bumpCounter(1, 3, 5)
	if runN != 4 || checkN != 5 {
		t.Errorf("bumpCounter(1, 3, 5) = (%d, %d), want (4, 5)", runN, checkN)
	}
	runN, checkN = bumpCounter(2, 3, 5)
	if runN != 3 || checkN != 6 {
		t.Errorf("bumpCounter(2, 3, 5) = (%d, %d), want (3, 6)", runN, checkN)
	}
}
