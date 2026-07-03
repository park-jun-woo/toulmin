//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestControlRoundTrip — SECTION/CASE/PROSE controls round-trip and empty PROSE Span is rejected
package packet

import (
	"reflect"
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// TestControlRoundTrip checks SECTION / CASE_BEGIN / CASE_END (subcode + Arg) and
// PROSE (subcode + SPAN) controls encode and decode identically, that the frame
// class is Control, and that a PROSE packet without a Span is rejected.
func TestControlRoundTrip(t *testing.T) {
	span := word.Span{TID: 1, Off: 0x0010, Len: 0x0020}
	cases := []ControlPacket{
		{Code: codebook.ControlSection, Arg: uint16(codebook.SectionCases)},
		{Code: codebook.ControlCaseBegin, Arg: 0x0010},
		{Code: codebook.ControlCaseEnd, Arg: 0x0010},
		{Code: codebook.ControlProse, Span: &span},
	}
	for _, in := range cases {
		ws, err := EncodeControl(in)
		if err != nil {
			t.Fatalf("EncodeControl(%v): %v", in, err)
		}
		f, _, err := word.DecodeFrame(ws)
		if err != nil || f.Class != word.ClassControl {
			t.Fatalf("frame class = %v, %v; want Control", f.Class, err)
		}
		got, n, err := DecodeControl(ws)
		if err != nil || n != len(ws) || !reflect.DeepEqual(got, in) {
			t.Fatalf("DecodeControl = %v, %d, %v; want %v", got, n, err, in)
		}
	}
	if _, err := EncodeControl(ControlPacket{Code: codebook.ControlProse}); err == nil {
		t.Fatalf("PROSE without Span: want error")
	}
}
