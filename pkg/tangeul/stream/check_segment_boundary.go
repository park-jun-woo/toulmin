//ff:func feature=tangl type=util control=sequence
//ff:what checkSegmentBoundary — verify the second packet is a SEGMENT_BOUNDARY control
package stream

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// checkSegmentBoundary verifies the packet after the header is a SEGMENT_BOUNDARY
// control (§7.1, one document = one segment). Its body is the subcode alone, so
// packet.DecodeControl (which expects subcode+Arg) does not model it; the check
// is done directly against the framed body here.
func checkSegmentBoundary(rp rawPacket) error {
	if rp.Class != word.ClassControl {
		return fmt.Errorf("tangeul/stream: expected SEGMENT_BOUNDARY control, got class %d", uint8(rp.Class))
	}
	_, head, err := word.DecodeFrame(rp.Words)
	if err != nil {
		return err
	}
	body := rp.Words[head:]
	if len(body) != 1 || codebook.ControlCode(body[0]) != codebook.ControlSegmentBoundary {
		return errors.New("tangeul/stream: malformed SEGMENT_BOUNDARY control")
	}
	return nil
}
