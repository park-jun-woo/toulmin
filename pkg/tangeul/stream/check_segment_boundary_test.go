//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCheckSegmentBoundary — a non-control class, a bad frame, a malformed body and a valid boundary
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestCheckSegmentBoundary covers checkSegmentBoundary's four outcomes: a
// non-Control frame class, a frame header word.DecodeFrame cannot parse, a
// well-framed body that is not exactly one SEGMENT_BOUNDARY control word, and
// the valid SEGMENT_BOUNDARY.
func TestCheckSegmentBoundary(t *testing.T) {
	validFrame := word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 1})

	cases := []struct {
		name    string
		rp      rawPacket
		wantErr bool
	}{
		{
			"non-control class",
			rawPacket{Class: word.ClassNode, Words: nil},
			true,
		},
		{
			"bad frame header",
			rawPacket{Class: word.ClassControl, Words: nil},
			true,
		},
		{
			"malformed body, wrong length",
			rawPacket{Class: word.ClassControl, Words: append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 2}), word.Word(codebook.ControlSegmentBoundary), 0x0000)},
			true,
		},
		{
			"malformed body, wrong control code",
			rawPacket{Class: word.ClassControl, Words: append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 1}), word.Word(codebook.ControlStreamHeader))},
			true,
		},
		{
			"valid segment boundary",
			rawPacket{Class: word.ClassControl, Words: append(append([]word.Word(nil), validFrame...), word.Word(codebook.ControlSegmentBoundary))},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkSegmentBoundary(c.rp)
			if c.wantErr && err == nil {
				t.Fatalf("checkSegmentBoundary(%+v): expected error, got nil", c.rp)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("checkSegmentBoundary(%+v): %v", c.rp, err)
			}
		})
	}
}
