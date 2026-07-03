//ff:func feature=tangl type=util control=selection
//ff:what emitStmt — dispatch a StmtBlock to its section emitter
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitStmt routes a statement to the emitter for its section. Each emitter
// returns the statement's packets with all non-span payloads first and the
// span-bearing Node/Edge run second, so the reader groups by that shape.
func emitStmt(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	switch ref.Sec {
	case codebook.SectionSubject:
		return emitSubject(e, ref, elem, span)
	case codebook.SectionSee:
		return emitSee(e, ref, elem, span)
	case codebook.SectionDefinitions:
		return emitDefinition(e, ref, elem, span)
	case codebook.SectionRules:
		return emitRule(e, ref, elem, span)
	case codebook.SectionCases:
		return emitCase(e, ref, elem, span)
	case codebook.SectionProvides:
		return emitProvides(e, ref, elem, span)
	case codebook.SectionInternal:
		return emitInternal(e, ref, elem, span)
	default:
		return nil, fmt.Errorf("tangeul/stream: unknown section %d", uint16(ref.Sec))
	}
}
