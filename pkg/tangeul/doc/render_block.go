//ff:func feature=tangl type=util control=selection
//ff:what renderBlock — render one Block for mode 2 (prose verbatim, section header, or codec statement)
package doc

import "fmt"

// renderBlock renders a single Block for mode 2 (README contract §3): a
// ProseBlock is replayed verbatim from Source, a SectionBlock heading is
// regenerated in the codec's locale, and a StmtBlock is rendered via the codec's
// RenderStmt with a core-owned list marker and section indent (top level, or one
// level for a body sub-item whose Path indexes into a sub-slice).
func renderBlock(d *Doc, c Codec, b Block, i int) ([]byte, error) {
	switch b.Kind {
	case ProseBlock:
		end := b.Span.Off + b.Span.Len
		if end < b.Span.Off || int(end) > len(d.Source) {
			return nil, fmt.Errorf("tangeul/doc: block %d span [%d,%d) out of Source (%d)", i, b.Span.Off, end, len(d.Source))
		}
		return d.Source[b.Span.Off:end], nil
	case SectionBlock:
		h, err := sectionHeader(c.Locale(), b.Ref.Sec)
		if err != nil {
			return nil, err
		}
		return append([]byte(h), '\n'), nil
	case StmtBlock:
		stmt, err := c.RenderStmt(d, b.Ref)
		if err != nil {
			return nil, err
		}
		level := 0
		if len(b.Ref.Path) >= 2 {
			level = 1
		}
		return renderStmtBlock(stmt, level), nil
	default:
		return nil, fmt.Errorf("tangeul/doc: block %d has unknown BlockKind %d", i, uint8(b.Kind))
	}
}
