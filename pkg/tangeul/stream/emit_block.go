//ff:func feature=tangl type=util control=selection
//ff:what emitBlock — encode one Block's packets (prose span, section heading, or statement)
package stream

import (
	"errors"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitBlock encodes a single Block: a ProseBlock is a PROSE control, a
// SectionBlock is a SECTION control plus a heading PROSE (whose span the reader
// lifts back onto the SectionBlock), and a StmtBlock resolves its AST element and
// delegates to emitStmt with the block span.
func emitBlock(e *enc, d *doc.Doc, b doc.Block) ([]word.Word, error) {
	switch b.Kind {
	case doc.ProseBlock:
		return proseWords(b.Span)
	case doc.SectionBlock:
		sec, err := caseControlWords(codebook.ControlSection, uint16(b.Ref.Sec))
		if err != nil {
			return nil, err
		}
		pr, err := proseWords(b.Span)
		if err != nil {
			return nil, err
		}
		return append(sec, pr...), nil
	case doc.StmtBlock:
		elem, err := doc.ResolveStmt(d, b.Ref)
		if err != nil {
			return nil, err
		}
		ws, err := toWireSpan(b.Span)
		if err != nil {
			return nil, err
		}
		return emitStmt(e, b.Ref, elem, ws)
	default:
		return nil, errors.New("tangeul/stream: unknown BlockKind in emitBlock")
	}
}
