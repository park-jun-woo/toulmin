//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeSpan — short input errors and a full 3-word Span round-trips
package word

import "testing"

func TestDecodeSpan(t *testing.T) {
	if _, err := DecodeSpan([]Word{1, 2}); err == nil {
		t.Fatal("expected error: SPAN needs 3 words")
	}
	want := Span{TID: 3, Off: 65000, Len: 128}
	ws := EncodeSpan(want)
	if got, err := DecodeSpan(ws[:]); err != nil || got != want {
		t.Fatalf("roundtrip: got %+v err=%v want %+v", got, err, want)
	}
}
