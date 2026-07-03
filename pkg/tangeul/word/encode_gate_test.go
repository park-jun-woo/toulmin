//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeGate — golden "at least 75%" = 0x004B, operator placement, percent>100 rejection
package word

import "testing"

func TestEncodeGate(t *testing.T) {
	// "if at least 75% certain" = op 0, percent 75 → 0x004B (§3.2).
	if w, err := EncodeGate(Gate{Op: GateAtLeast, Percent: 75}); err != nil || w != 0x004B {
		t.Fatalf("at least 75%%: got 0x%04X err=%v want 0x004B", w, err)
	}
	// Operator lands in bits 15-14.
	if w, err := EncodeGate(Gate{Op: GateAtMost, Percent: 100}); err != nil || w != Word(0x3<<14|100) {
		t.Fatalf("at most 100%%: got 0x%04X err=%v", w, err)
	}
	// percent over 100 is rejected.
	if _, err := EncodeGate(Gate{Op: GateAbove, Percent: 101}); err == nil {
		t.Fatal("expected error for percent 101")
	}
}
