//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestTextRoundTrip — Text payloads round-trip for even/odd/empty/multibyte and reject a bad format
package packet

import (
	"bytes"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestTextRoundTrip checks Text payloads round-trip for even, odd, empty and
// multibyte UTF-8 byte lengths, that the format word and NUM32 byte length are
// emitted, that odd lengths pad with a single trailing 0x00, and that a wrong
// format code is rejected.
func TestTextRoundTrip(t *testing.T) {
	cases := [][]byte{
		nil,
		[]byte("A"),
		[]byte("AB"),
		[]byte("in case of `can withdraw`\n"),
		[]byte("금액이 충분한가"), // multibyte UTF-8, verbatim bytes
	}
	for _, b := range cases {
		ws := EncodeText(TextPacket{Bytes: b})
		if ws[1] != word.Word(TextFormatUTF8) {
			t.Fatalf("format word = %04X, want %04X", ws[1], uint16(TextFormatUTF8))
		}
		bl, _ := word.DecodeNum32(ws[2:4])
		if int(bl) != len(b) {
			t.Fatalf("byte length = %d, want %d", bl, len(b))
		}
		got, n, err := DecodeText(ws)
		if err != nil || n != len(ws) || !bytes.Equal(got.Bytes, b) {
			t.Fatalf("DecodeText(%q) = %q, %d, %v", b, got.Bytes, n, err)
		}
	}
	// A corrupted format code must be rejected, not guessed.
	bad := EncodeText(TextPacket{Bytes: []byte("x")})
	bad[1] = 0x00FF
	if _, _, err := DecodeText(bad); err == nil {
		t.Fatalf("bad format code: want error")
	}
}
