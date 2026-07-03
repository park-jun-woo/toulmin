//ff:func feature=tangl type=util control=selection
//ff:what EncodeControl — encode a ControlPacket (subcode + Arg word, or subcode + SPAN for PROSE)
package packet

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// EncodeControl encodes p as a Control packet (§2.4/§7). PROSE carries the
// subcode then a 3-word SPAN (§7.4); every other supported control (SECTION,
// CASE_BEGIN, CASE_END) carries the subcode then a single Arg word (SectionCode
// or Case TID). A PROSE packet without a Span is a hard error.
func EncodeControl(p ControlPacket) ([]word.Word, error) {
	body := []word.Word{word.Word(p.Code)}
	switch p.Code {
	case codebook.ControlProse:
		if p.Span == nil {
			return nil, errors.New("tangeul/packet: PROSE control requires a Span")
		}
		s := word.EncodeSpan(*p.Span)
		body = append(body, s[0], s[1], s[2])
	default:
		body = append(body, word.Word(p.Arg))
	}
	return append(word.EncodeFrame(word.Frame{Class: word.ClassControl, Len: uint32(len(body))}), body...), nil
}
