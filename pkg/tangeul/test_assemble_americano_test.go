//ff:func feature=tangl type=util control=sequence
//ff:what assembleAmericanoDoc — flat Block tiling for americano.md paired with its parsed AST
package tangeul

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// assembleAmericanoDoc lays out the Blocks for americano.md: Subject, two See
// aliases (arm, sensor), one const Definition, four brew-chain cases, a Provides
// endpoint, and an Internal `every` trigger. It exists only so TestAmericano
// Excluded can drive Save far enough to hit the `do log.CoffeeStatus` exec, whose
// `log` alias is never `see`n — the stream Writer rejects it, which is why
// americano has no golden. Every block gets a 4-byte placeholder span tiling the
// Source exactly.
func assembleAmericanoDoc(a *ast.Document) *doc.Doc {
	var src []byte
	var blocks []doc.Block
	add := func(kind doc.BlockKind, ref doc.StmtRef) {
		off := uint32(len(src))
		src = append(src, '-', '-', '-', '-')
		blocks = append(blocks, doc.Block{Kind: kind, Span: doc.Span{Off: off, Len: 4}, Ref: ref})
	}
	sr := func(sec codebook.SectionCode, p ...int) doc.StmtRef { return doc.StmtRef{Sec: sec, Path: p} }

	add(doc.SectionBlock, sr(codebook.SectionSubject))
	add(doc.StmtBlock, sr(codebook.SectionSubject))
	add(doc.SectionBlock, sr(codebook.SectionSee))
	add(doc.StmtBlock, sr(codebook.SectionSee, 0))
	add(doc.StmtBlock, sr(codebook.SectionSee, 1))
	add(doc.SectionBlock, sr(codebook.SectionDefinitions))
	add(doc.StmtBlock, sr(codebook.SectionDefinitions, 0))
	add(doc.SectionBlock, sr(codebook.SectionCases))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 2, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 3, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 2))
	add(doc.StmtBlock, sr(codebook.SectionCases, 2, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 2, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 2, 3, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3, 1, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3, 2, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 3, 3, 1))
	add(doc.SectionBlock, sr(codebook.SectionProvides))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 1, 0))
	add(doc.SectionBlock, sr(codebook.SectionInternal))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 0))
	add(doc.StmtBlock, sr(codebook.SectionInternal, 0, 1, 0))

	return &doc.Doc{AST: a, Source: src, Blocks: blocks, Locale: "en"}
}
