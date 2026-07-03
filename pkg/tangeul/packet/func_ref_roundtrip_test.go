//ff:func feature=tangl type=util control=sequence
//ff:what TestFuncRefRoundTrip — FuncRef encodes to its two words and decodes back (Label excluded)
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestFuncRefRoundTrip checks that encodeFuncRef emits exactly the Alias and
// Verb words and decodeFuncRef recovers them (Label is stream-serialized, so it
// is not part of the encoded round-trip).
func TestFuncRefRoundTrip(t *testing.T) {
	in := FuncRef{Alias: 0x0002, Verb: 0xC200}
	ws := encodeFuncRef(in)
	if len(ws) != 2 || ws[0] != word.Word(0x0002) || ws[1] != word.Word(0xC200) {
		t.Fatalf("encodeFuncRef = %v, want [0002 C200]", ws)
	}
	got, err := decodeFuncRef(ws)
	if err != nil || got != in {
		t.Fatalf("decodeFuncRef = %v, %v; want %v", got, err, in)
	}
	if _, err := decodeFuncRef(ws[:1]); err == nil {
		t.Fatalf("decodeFuncRef of short slice: want error")
	}
}
