//ff:type feature=tangl type=model
//ff:what BlockKind — the tiling role of a Block (prose, statement, or section header)
package doc

// BlockKind classifies a Block's role in the Source Text tiling.
type BlockKind uint8

const (
	// ProseBlock covers authored prose bytes replayed verbatim in both modes.
	ProseBlock BlockKind = iota
	// StmtBlock anchors an AST statement element; mode 2 regenerates it via a Codec.
	StmtBlock
	// SectionBlock covers a `## tangl:<name>` heading; mode 2 regenerates it per locale.
	SectionBlock
)
