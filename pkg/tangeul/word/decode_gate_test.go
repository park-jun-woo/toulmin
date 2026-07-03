//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeGate — golden 0x004B decode, reserved-bit and percent-range violations are errors
package word

import "testing"

func TestDecodeGate(t *testing.T) {
	// Golden decode.
	if g, err := DecodeGate(0x004B); err != nil || g.Op != GateAtLeast || g.Percent != 75 {
		t.Fatalf("0x004B: got %+v err=%v", g, err)
	}
	// Any reserved bit (13-7) set is a hard error.
	if _, err := DecodeGate(0x0080); err == nil {
		t.Fatal("expected reserved-bit error")
	}
	// percent field 127 with clean reserved bits still violates the 0-100 range.
	if _, err := DecodeGate(0x007F); err == nil {
		t.Fatal("expected percent-range error")
	}
}
