//ff:func feature=tangl type=util control=sequence
//ff:what TestProseWords — a valid span encodes a PROSE control, and a span exceeding the 16-bit wire width errors
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// TestProseWords covers both branches of proseWords: toWireSpan succeeding
// (the PROSE control round-trips) and toWireSpan failing on a span that
// overflows the 16-bit wire width, propagated without wrapping in a Control
// frame.
func TestProseWords(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		got, err := proseWords(doc.Span{Off: 10, Len: 20})
		if err != nil {
			t.Fatalf("proseWords: %v", err)
		}
		cp, n, err := packet.DecodeControl(got)
		if err != nil {
			t.Fatalf("DecodeControl: %v", err)
		}
		if n != len(got) {
			t.Errorf("consumed %d words, want %d", n, len(got))
		}
		if cp.Code != codebook.ControlProse || cp.Span == nil {
			t.Fatalf("got Code=%v Span=%v, want ControlProse with a Span", cp.Code, cp.Span)
		}
		if cp.Span.Off != 10 || cp.Span.Len != 20 {
			t.Errorf("Span = %+v, want Off=10 Len=20", cp.Span)
		}
	})

	t.Run("spanOverflow", func(t *testing.T) {
		_, err := proseWords(doc.Span{Off: 0x10000, Len: 1})
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})
}
