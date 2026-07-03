//ff:func feature=tangl type=parser control=sequence
//ff:what TestRenderRunCheck — renderRunCheck maps slotKind 1/2 to run/check and rejects others
package en

import "testing"

// TestRenderRunCheck checks the run (slotKind 1) and check (slotKind 2) forms and
// that an unknown slotKind is an error.
func TestRenderRunCheck(t *testing.T) {
	if got, err := renderRunCheck(1, "can go"); err != nil || got != "run `can go`" {
		t.Errorf("run: got %q, err %v", got, err)
	}
	if got, err := renderRunCheck(2, "can sub"); err != nil || got != "check `can sub`" {
		t.Errorf("check: got %q, err %v", got, err)
	}
	if _, err := renderRunCheck(0, "x"); err == nil {
		t.Errorf("slotKind 0: expected error")
	}
}
