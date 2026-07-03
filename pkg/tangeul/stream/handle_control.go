//ff:func feature=tangl type=util control=selection
//ff:what handleControl — apply a framing control (SECTION, PROSE, CASE_BEGIN/END) to the reader state
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// handleControl decodes a framing control and updates the reader: SECTION arms
// the pending section (its heading PROSE follows), PROSE becomes a SectionBlock
// (when a section is pending) or a ProseBlock, CASE_END clears the open case, and
// CASE_BEGIN is a no-op (membership is set when the case node is read).
func handleControl(d *dec, pendingSection *int, rp rawPacket) error {
	cp, _, err := packet.DecodeControl(rp.Words)
	if err != nil {
		return err
	}
	if *pendingSection >= 0 && cp.Code != codebook.ControlProse {
		return fmt.Errorf("tangeul/stream: SECTION control must be followed by its heading PROSE")
	}
	switch cp.Code {
	case codebook.ControlSection:
		*pendingSection = int(cp.Arg)
	case codebook.ControlProse:
		addSpanBlock(d, pendingSection, fromWireSpan(*cp.Span))
	case codebook.ControlCaseBegin:
		// membership is set when the case node is read.
	case codebook.ControlCaseEnd:
		d.curCase = -1
	default:
		return fmt.Errorf("tangeul/stream: unexpected control 0x%02X mid-stream", uint16(cp.Code))
	}
	return nil
}
