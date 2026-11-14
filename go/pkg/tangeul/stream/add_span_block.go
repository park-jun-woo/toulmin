//ff:func feature=tangl type=util control=sequence
//ff:what addSpanBlock — record a PROSE span as a SectionBlock heading (pending SECTION) or a ProseBlock
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// addSpanBlock appends the Block a PROSE record describes: the heading of the
// pending SECTION (clearing it) or a free-standing ProseBlock.
func addSpanBlock(d *dec, pendingSection *int, ws toulmin.Span) error {
	span, err := fromWireSpan(ws)
	if err != nil {
		return err
	}
	if *pendingSection >= 0 {
		d.blocks = append(d.blocks, doc.Block{Kind: doc.SectionBlock, Span: span, Ref: doc.StmtRef{Sec: toulmin.SectionCode(*pendingSection)}})
		*pendingSection = -1
		return nil
	}
	d.blocks = append(d.blocks, doc.Block{Kind: doc.ProseBlock, Span: span})
	return nil
}
