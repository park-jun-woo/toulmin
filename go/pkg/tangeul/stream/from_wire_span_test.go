//ff:func feature=tangl type=util control=sequence
//ff:what TestFromWireSpan — fromWireSpan widens a wire toulmin.Span's Off/Len into a doc.Span, dropping the TID
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestFromWireSpan checks the straight-line widening: Off and Len carry over
// and a value past 32 bits is reported, not truncated.
func TestFromWireSpan(t *testing.T) {
	got, err := fromWireSpan(toulmin.Span{Off: 12, Len: 34})
	want := doc.Span{Off: 12, Len: 34}
	if err != nil || got != want {
		t.Errorf("fromWireSpan = %+v %v, want %+v", got, err, want)
	}
	if _, err := fromWireSpan(toulmin.Span{Off: 1 << 32}); err == nil {
		t.Error("offset past 32 bits: want error")
	}
	if _, err := fromWireSpan(toulmin.Span{Len: 1 << 32}); err == nil {
		t.Error("length past 32 bits: want error")
	}
}
