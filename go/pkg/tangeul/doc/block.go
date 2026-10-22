//ff:type feature=tangl type=model
//ff:what Block — one stream-order tile: a Source span plus its kind and AST coordinate
package doc

// Block is one tile of the Source Text in stream (document) order. Span is the
// covered byte range [Off, Off+Len) into Doc.Source; Kind classifies the tile;
// Ref is the AST element coordinate when Kind is StmtBlock (and carries the
// SectionCode for a SectionBlock header). The stream order of Blocks is the
// document order — the single source of truth for element ordering, since
// ast.Case stores Requires/Nodes/Attacks/Execs in separate slices.
type Block struct {
	Kind BlockKind
	Span Span
	Ref  StmtRef
}
