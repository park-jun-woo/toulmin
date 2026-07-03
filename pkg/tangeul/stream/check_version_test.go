//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCheckVersion — non-control class, bad frame, short body, wrong control, and each version-major mismatch
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestCheckVersion covers checkVersion's outcomes: a non-Control frame class,
// a frame header word.DecodeFrame cannot parse, a body shorter than 5 words, a
// first control word that is not STREAM_HEADER, a GEUL major mismatch, a
// TANGEUL profile mismatch, a codebook major mismatch, and a fully valid stamp.
func TestCheckVersion(t *testing.T) {
	// makeVersionBody builds a well-formed STREAM_HEADER body (5 words) with
	// the given majors/profile so each mismatch check in checkVersion can be
	// exercised independently.
	makeVersionBody := func(control codebook.ControlCode, geulMajor uint8, profile uint16, cbMajor uint8) []word.Word {
		return []word.Word{
			word.Word(control),
			word.Word(uint16(geulMajor)<<8 | uint16(codebook.GEULVersionMinor)),
			word.Word(uint16(codebook.GEULVersionPatch) << 8),
			word.Word(profile),
			word.Word(uint16(cbMajor)<<8 | uint16(codebook.CodebookVersionMinor)),
		}
	}
	frame := func(n int) []word.Word { return word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(n)}) }
	withFrame := func(body []word.Word) []word.Word {
		return append(append([]word.Word(nil), frame(len(body))...), body...)
	}

	cases := []struct {
		name    string
		rp      rawPacket
		wantErr bool
	}{
		{"non-control class", rawPacket{Class: word.ClassNode, Words: nil}, true},
		{"bad frame header", rawPacket{Class: word.ClassControl, Words: nil}, true},
		{
			"body too short",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlStreamHeader, codebook.GEULVersionMajor, codebook.TANGLProfileID, codebook.CodebookVersionMajor)[:4])},
			true,
		},
		{
			"wrong first control",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlSegmentBoundary, codebook.GEULVersionMajor, codebook.TANGLProfileID, codebook.CodebookVersionMajor))},
			true,
		},
		{
			"geul major mismatch",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlStreamHeader, codebook.GEULVersionMajor+1, codebook.TANGLProfileID, codebook.CodebookVersionMajor))},
			true,
		},
		{
			"tangl profile mismatch",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlStreamHeader, codebook.GEULVersionMajor, codebook.TANGLProfileID+1, codebook.CodebookVersionMajor))},
			true,
		},
		{
			"codebook major mismatch",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlStreamHeader, codebook.GEULVersionMajor, codebook.TANGLProfileID, codebook.CodebookVersionMajor+1))},
			true,
		},
		{
			"valid version stamp",
			rawPacket{Class: word.ClassControl, Words: withFrame(makeVersionBody(codebook.ControlStreamHeader, codebook.GEULVersionMajor, codebook.TANGLProfileID, codebook.CodebookVersionMajor))},
			false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := checkVersion(c.rp)
			if c.wantErr && err == nil {
				t.Fatalf("checkVersion(%+v): expected error, got nil", c.rp)
			}
			if !c.wantErr && err != nil {
				t.Fatalf("checkVersion(%+v): %v", c.rp, err)
			}
		})
	}
}
