//ff:func feature=tangl type=util control=sequence
//ff:what TestEmitHeader — decodes the two leading frames and checks their Class/Len/subcode shape
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestEmitHeader is a structural companion to TestHeaderGolden (which pins the
// exact hex bytes): it decodes emitHeader's output via word.DecodeFrame and
// checks the STREAM_HEADER frame is a 6-word Control carrying
// ControlStreamHeader, and the following SEGMENT_BOUNDARY frame is a 1-word
// Control carrying ControlSegmentBoundary — the shape a reader actually relies
// on. emitHeader itself is a straight-line sequence with no branches, so a
// single call reaches full statement coverage.
func TestEmitHeader(t *testing.T) {
	got := emitHeader(codebook.LocaleUnknown)

	f1, n1, err := word.DecodeFrame(got)
	if err != nil {
		t.Fatalf("DecodeFrame(header): %v", err)
	}
	if f1.Class != word.ClassControl || f1.Len != 6 {
		t.Fatalf("header frame = %+v, want Class=Control Len=6", f1)
	}
	if got[n1] != word.Word(codebook.ControlStreamHeader) {
		t.Errorf("subcode = %#x, want ControlStreamHeader", got[n1])
	}

	rest := got[n1+6:]
	f2, n2, err := word.DecodeFrame(rest)
	if err != nil {
		t.Fatalf("DecodeFrame(boundary): %v", err)
	}
	if f2.Class != word.ClassControl || f2.Len != 1 {
		t.Fatalf("boundary frame = %+v, want Class=Control Len=1", f2)
	}
	if rest[n2] != word.Word(codebook.ControlSegmentBoundary) {
		t.Errorf("subcode = %#x, want ControlSegmentBoundary", rest[n2])
	}
	if len(got) != n1+6+n2+1 {
		t.Errorf("emitHeader() len = %d, want %d", len(got), n1+6+n2+1)
	}
}
