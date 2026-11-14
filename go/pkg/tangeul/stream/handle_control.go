//ff:func feature=tangl type=util control=selection
//ff:what handleControl — apply a framing record: arm a SECTION, place PROSE, open or close a case
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// handleControl decodes a framing record and updates the reader: SECTION arms
// the pending section (its heading PROSE follows), PROSE becomes a SectionBlock
// (when a section is pending) or a ProseBlock, CASE-END clears the open case, and
// CASE-BEGIN is a no-op (membership is set when the case node is read). The
// record's Edge Id is registered like every other record's.
func handleControl(d *dec, pendingSection *int, rp rawPacket) error {
	cp, err := toulmin.DecodeControl(rp.Rec, d.wire)
	if err != nil {
		return err
	}
	if err := claimID(d, cp.ID); err != nil {
		return err
	}
	if *pendingSection >= 0 && cp.Code != toulmin.ControlProse {
		return fmt.Errorf("tangeul/stream: SECTION control must be followed by its heading PROSE")
	}
	switch cp.Code {
	case toulmin.ControlSection:
		*pendingSection = int(cp.Arg)
	case toulmin.ControlProse:
		return addSpanBlock(d, pendingSection, *cp.Span)
	case toulmin.ControlCaseEnd:
		d.curCase = -1
	}
	return nil
}
