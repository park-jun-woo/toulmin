//ff:func feature=tangl type=util control=sequence
//ff:what caseControlWords — encode a CASE_BEGIN/END or SECTION control carrying a single arg word
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// caseControlWords encodes a Control packet whose body is a single argument word:
// CASE_BEGIN/CASE_END carry a Case TID, SECTION carries a SectionCode.
func caseControlWords(code codebook.ControlCode, arg uint16) ([]word.Word, error) {
	return packet.EncodeControl(packet.ControlPacket{Code: code, Arg: arg})
}
