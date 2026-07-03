//ff:func feature=tangl type=util control=sequence
//ff:what TestDecodeControlErrors — DecodeControl rejects a bad frame, truncated bodies and PROSE/default body-length mismatches
package packet

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestDecodeControlErrors exercises DecodeControl's self-contained guards: a
// malformed frame header, a body shorter than the framed length, a body too
// short for a subcode, a PROSE body that is not subcode+SPAN, and a
// non-PROSE body that is not subcode+Arg.
func TestDecodeControlErrors(t *testing.T) {
	if _, _, err := DecodeControl([]word.Word{0x0000}); err == nil {
		t.Fatalf("bad frame prefix: want error")
	}

	full, err := EncodeControl(ControlPacket{Code: codebook.ControlCaseBegin, Arg: 1})
	if err != nil {
		t.Fatalf("EncodeControl: %v", err)
	}
	if _, _, err := DecodeControl(full[:len(full)-1]); err == nil {
		t.Fatalf("truncated body: want error")
	}

	empty := word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 0})
	if _, _, err := DecodeControl(empty); err == nil {
		t.Fatalf("body shorter than a subcode: want error")
	}

	proseShort := append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 2}), word.Word(codebook.ControlProse), word.Word(0))
	if _, _, err := DecodeControl(proseShort); err == nil {
		t.Fatalf("PROSE body not subcode+SPAN: want error")
	}

	defaultShort := append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: 1}), word.Word(codebook.ControlCaseBegin))
	if _, _, err := DecodeControl(defaultShort); err == nil {
		t.Fatalf("default body not subcode+Arg: want error")
	}
}
