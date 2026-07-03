//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeSpan — a Span encodes to 3 words in TID/offset/length order
package word

import "testing"

func TestEncodeSpan(t *testing.T) {
	got := EncodeSpan(Span{TID: 0x0011, Off: 0x1234, Len: 0x0040})
	if got != [spanWords]Word{0x0011, 0x1234, 0x0040} {
		t.Fatalf("got %v", got)
	}
}
