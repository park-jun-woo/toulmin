//ff:func feature=tangl type=util control=sequence
//ff:what readPreamble — require the SOURCE record first in the BODY and read its text and SPAN width
package stream

import (
	"errors"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// readPreamble requires that the BODY opens with the SOURCE record (GEUL t002
// TOULMIN 문서) and returns its text with the SPAN width it declares for every
// later record.
func readPreamble(packets []rawPacket) (toulmin.TextPacket, int, error) {
	if len(packets) == 0 || packets[0].Class != classSource {
		return toulmin.TextPacket{}, 0, errors.New("tangeul/stream: BODY must open with the SOURCE record")
	}
	return toulmin.DecodeSource(packets[0].Rec)
}
