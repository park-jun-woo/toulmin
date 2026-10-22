//ff:func feature=tangl type=util control=sequence
//ff:what TestProseWords — proseWords encodes a PROSE record and rejects a span the width cannot hold
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// TestProseWords checks the PROSE record carries the span, and that a span past
// the declared 1-word SPAN width is an error rather than a truncation.
func TestProseWords(t *testing.T) {
	got, err := proseWords(testEnc(nil), doc.Span{Off: 10, Len: 20})
	if err != nil {
		t.Fatalf("proseWords: %v", err)
	}
	cps := controlsOf(t, got)
	if len(cps) != 1 || cps[0].Code != toulmin.ControlProse || cps[0].Span == nil || cps[0].Span.Off != 10 || cps[0].Span.Len != 20 {
		t.Fatalf("got %+v, want one PROSE 10+20", cps)
	}
	if _, err := proseWords(testEnc(nil), doc.Span{Off: 0x10000, Len: 1}); err == nil {
		t.Fatal("span past the SPAN width: want error")
	}
}
