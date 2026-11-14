//ff:func feature=tangl type=util control=selection
//ff:what emitStmt — dispatch a StmtBlock to its section emitter
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// emitStmt routes a statement to the emitter for its section. Each emitter
// returns the statement's packets with all non-span payloads first and the
// span-bearing Node/Edge run second, so the reader groups by that shape.
func emitStmt(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	switch ref.Sec {
	case toulmin.SectionSubject:
		return emitSubject(e, ref, elem, span)
	case toulmin.SectionSee:
		return emitSee(e, ref, elem, span)
	case toulmin.SectionDefinitions:
		return emitDefinition(e, ref, elem, span)
	case toulmin.SectionRules:
		return emitRule(e, ref, elem, span)
	case toulmin.SectionCases:
		return emitCase(e, ref, elem, span)
	case toulmin.SectionProvides:
		return emitProvides(e, ref, elem, span)
	case toulmin.SectionInternal:
		return emitInternal(e, ref, elem, span)
	default:
		return nil, fmt.Errorf("tangeul/stream: unknown section %d", uint16(ref.Sec))
	}
}
