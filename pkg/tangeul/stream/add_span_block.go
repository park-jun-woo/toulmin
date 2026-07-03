//ff:func feature=tangl type=util control=sequence
//ff:what addSpanBlock — append a SectionBlock (consuming a pending section) or a ProseBlock
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
)

// addSpanBlock records a PROSE control's span as a Block: a SectionBlock when a
// SECTION is pending (its heading), otherwise a standalone ProseBlock. The
// pending section is cleared once its heading is consumed.
func addSpanBlock(d *dec, pendingSection *int, span doc.Span) {
	if *pendingSection >= 0 {
		d.blocks = append(d.blocks, doc.Block{Kind: doc.SectionBlock, Span: span, Ref: doc.StmtRef{Sec: codebook.SectionCode(*pendingSection)}})
		*pendingSection = -1
		return
	}
	d.blocks = append(d.blocks, doc.Block{Kind: doc.ProseBlock, Span: span})
}
