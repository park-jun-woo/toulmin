//ff:func feature=tangl type=util control=sequence
//ff:what TestFromWireSpan — fromWireSpan widens a wire word.Span's Off/Len into a doc.Span, dropping the TID
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestFromWireSpan checks the straight-line widening: Off and Len carry over
// and the anchor TID is dropped.
func TestFromWireSpan(t *testing.T) {
	got := fromWireSpan(word.Span{TID: 9, Off: 12, Len: 34})
	want := doc.Span{Off: 12, Len: 34}
	if got != want {
		t.Errorf("fromWireSpan = %+v, want %+v", got, want)
	}
}
