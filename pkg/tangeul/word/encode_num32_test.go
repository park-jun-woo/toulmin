//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeNum32 — a uint32 encodes to two words, high word first
package word

import "testing"

func TestEncodeNum32(t *testing.T) {
	if got := EncodeNum32(0x12345678); got[0] != 0x1234 || got[1] != 0x5678 {
		t.Fatalf("got [0x%04X 0x%04X] want [0x1234 0x5678]", got[0], got[1])
	}
}
