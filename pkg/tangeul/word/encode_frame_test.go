//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeFrame — golden checks for single-word headers and the 1022/1023 body-length boundary
package word

import "testing"

func TestEncodeFrame(t *testing.T) {
	// Node class, small body: [0001 | 00 | Len].
	if got := EncodeFrame(Frame{Class: ClassNode, Len: 5}); len(got) != 1 || got[0] != Word(0x1<<12|0x0<<10|5) {
		t.Fatalf("Node Len=5: got %v", got)
	}
	// Control class carries into bits 11-10.
	if got := EncodeFrame(Frame{Class: ClassControl, Len: 3}); got[0] != Word(0x1<<12|0x3<<10|3) {
		t.Fatalf("Control Len=3: got 0x%04X", got[0])
	}
	// 1022 is the largest body count that still fits a single header word.
	if got := EncodeFrame(Frame{Class: ClassPayload, Len: 1022}); len(got) != 1 || got[0]&0x3FF != 1022 {
		t.Fatalf("Len=1022 must be single word: got %v", got)
	}
	// 1023 triggers the 32-bit extension: sentinel Len + NUM32(real count).
	if got := EncodeFrame(Frame{Class: ClassPayload, Len: 1023}); len(got) != 3 || got[0]&0x3FF != frameLenSentinel || got[1] != 0 || got[2] != 1023 {
		t.Fatalf("Len=1023 extension: got %v", got)
	}
}
