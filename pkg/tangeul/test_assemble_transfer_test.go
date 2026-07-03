//ff:func feature=tangl type=util control=sequence
//ff:what assembleTransferDoc — flat Block tiling for transfer.md paired with its parsed AST
package tangeul

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// assembleTransferDoc lays out the Blocks for transfer.md in declaration order:
// Subject, two See aliases, a Cases section with `can withdraw` (2 requires, 2
// rule nodes, 1 don't-attack, do/undo/run execs) and `can deposit` (1 require, 1
// node, do + do-log execs), and a Provides endpoint (3 requires, 1 run). Each
// block gets a 4-byte placeholder span that tiles the Source exactly. The parsed
// AST a is carried unchanged so the reloaded stream matches the text-parse AST.
func assembleTransferDoc(a *ast.Document) *doc.Doc {
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
	add(doc.SectionBlock, sr(codebook.SectionCases))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 0, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 0, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 1, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 2, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 0, 3, 2))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 0, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 1, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 3, 0))
	add(doc.StmtBlock, sr(codebook.SectionCases, 1, 3, 1))
	add(doc.SectionBlock, sr(codebook.SectionProvides))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 0, 0))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 0, 1))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 0, 2))
	add(doc.StmtBlock, sr(codebook.SectionProvides, 0, 1, 0))

	return &doc.Doc{AST: a, Source: src, Blocks: blocks, Locale: "en"}
}
