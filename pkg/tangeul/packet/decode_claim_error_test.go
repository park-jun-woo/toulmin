//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeClaimErrors — decodeClaim rejects a short slice, a truncated object slot, an unknown ObjKind and a missing mode word
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeClaimErrors exercises decodeClaim's guards: fewer than 2 words
// (no subject/predicate), a truncated ObjTID/ObjNum/ObjSpan object slot, an
// unknown ObjKind, and a body that ends before the trailing mode word.
func TestDecodeClaimErrors(t *testing.T) {
	if _, _, err := decodeClaim([]word.Word{1}, ObjTID); err == nil {
		t.Fatalf("fewer than 2 words: want error")
	}
	if _, _, err := decodeClaim([]word.Word{1, 2}, ObjTID); err == nil {
		t.Fatalf("truncated ObjTID slot: want error")
	}
	if _, _, err := decodeClaim([]word.Word{1, 2}, ObjNum); err == nil {
		t.Fatalf("truncated ObjNum slot: want error")
	}
	if _, _, err := decodeClaim([]word.Word{1, 2}, ObjSpan); err == nil {
		t.Fatalf("truncated ObjSpan slot: want error")
	}
	if _, _, err := decodeClaim([]word.Word{1, 2}, ObjKind(3)); err == nil {
		t.Fatalf("unknown ObjKind: want error")
	}
	if _, _, err := decodeClaim([]word.Word{1, 2, 3}, ObjTID); err == nil {
		t.Fatalf("missing mode word: want error")
	}
}
