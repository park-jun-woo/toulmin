//ff:func feature=tangl type=util control=selection
//ff:what emitBlock — encode one Block's packets (prose span, section heading, or statement)
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// emitBlock encodes a single Block: a ProseBlock is a PROSE record, a
// SectionBlock is a SECTION control plus a heading PROSE (whose span the reader
// lifts back onto the SectionBlock), and a StmtBlock resolves its AST element and
// delegates to emitStmt with the block span.
func emitBlock(e *enc, d *doc.Doc, b doc.Block) ([]core.Word, error) {
	switch b.Kind {
	case doc.ProseBlock:
		return proseWords(e, b.Span), nil
	case doc.SectionBlock:
		return append(controlWords(e, toulmin.ControlSection, uint64(b.Ref.Sec), nil), proseWords(e, b.Span)...), nil
	case doc.StmtBlock:
		elem, err := doc.ResolveStmt(d, b.Ref)
		if err != nil {
			return nil, err
		}
		return emitStmt(e, b.Ref, elem, toWireSpan(b.Span))
	default:
		return nil, errors.New("tangeul/stream: unknown BlockKind in emitBlock")
	}
}
