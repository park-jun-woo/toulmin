//ff:func feature=tangl type=util control=sequence
//ff:what TestEncodeLabels — encodeLabels serializes an empty table, an even-length label and an odd-length label
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEncodeLabels checks encodeLabels handles an empty table, an
// even-length label (no padding) and an odd-length label (single trailing
// 0x00 pad).
func TestEncodeLabels(t *testing.T) {
	if got := encodeLabels(nil); len(got) != 1 || got[0] != word.Word(0) {
		t.Fatalf("empty table: got %v", got)
	}

	got := encodeLabels([]string{"ab", "c"})
	if got[0] != word.Word(2) {
		t.Fatalf("count word: got %v, want 2", got[0])
	}
	// "ab": byteLen=2 (even, no pad) -> count(1) + NUM32(2) + data(1) = 4 words.
	// "c": byteLen=1 (odd, one pad byte) -> NUM32(2) + data(1) = 3 words.
	if len(got) != 1+3+3 {
		t.Fatalf("total words = %d, want %d: %v", len(got), 1+3+3, got)
	}
}
