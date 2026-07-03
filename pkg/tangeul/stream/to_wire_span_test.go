//ff:func feature=tangl type=util control=sequence
//ff:what TestToWireSpan — toWireSpan accepts in-range spans and rejects each of the three 16-bit overflow shapes
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestToWireSpan checks the success path (offset/length both fit in 16 bits and
// their sum does too) and each way the guard clause can fail: Off alone over
// 0xFFFF, Len alone over 0xFFFF, and both individually in range but their sum
// overflowing. Each case asserts on a different branch of the single `if`, so
// the test itself selects distinct setup/assertions per case.
func TestToWireSpan(t *testing.T) {
	t.Run("in range", func(t *testing.T) {
		got, err := toWireSpan(doc.Span{Off: 10, Len: 20})
		if err != nil {
			t.Fatalf("toWireSpan: unexpected error: %v", err)
		}
		want := word.Span{TID: srcAnchorTID, Off: 10, Len: 20}
		if got != want {
			t.Errorf("toWireSpan = %+v, want %+v", got, want)
		}
	})

	t.Run("off overflow", func(t *testing.T) {
		_, err := toWireSpan(doc.Span{Off: 0x10000, Len: 0})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("len overflow", func(t *testing.T) {
		_, err := toWireSpan(doc.Span{Off: 0, Len: 0x10000})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("sum overflow", func(t *testing.T) {
		_, err := toWireSpan(doc.Span{Off: 0xFFF0, Len: 0x20})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("boundary exact fit", func(t *testing.T) {
		got, err := toWireSpan(doc.Span{Off: 0xFFF0, Len: 0xF})
		if err != nil {
			t.Fatalf("toWireSpan: unexpected error: %v", err)
		}
		want := word.Span{TID: srcAnchorTID, Off: 0xFFF0, Len: 0xF}
		if got != want {
			t.Errorf("toWireSpan = %+v, want %+v", got, want)
		}
	})
}
