//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeNum32 — a single-word input errors (NUM32 requires two words)
package word

import "testing"

func TestDecodeNum32(t *testing.T) {
	if _, err := DecodeNum32([]Word{1}); err == nil {
		t.Fatal("expected error: NUM32 needs 2 words")
	}
}
