//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeFrame — prefix violation, single-word decode, and truncated-extension error checks
package word

import "testing"

func TestDecodeFrame(t *testing.T) {
	// Prefix 0000 is not the Free prefix (0001).
	if _, _, err := DecodeFrame([]Word{0x0005}); err == nil {
		t.Fatal("expected prefix mismatch error, got nil")
	}
	// Empty input is an error.
	if _, _, err := DecodeFrame(nil); err == nil {
		t.Fatal("expected error on empty input")
	}
	// Single-word Edge header.
	f, n, err := DecodeFrame([]Word{0x1<<12 | 0x1<<10 | 42})
	if err != nil || n != 1 || f.Class != ClassEdge || f.Len != 42 {
		t.Fatalf("single word: f=%+v n=%d err=%v", f, n, err)
	}
	// Sentinel header without the trailing NUM32 is an error.
	if _, _, err := DecodeFrame([]Word{0x1<<12 | frameLenSentinel}); err == nil {
		t.Fatal("expected truncated-extension error")
	}
}
