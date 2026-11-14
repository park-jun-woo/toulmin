//ff:func feature=tangl type=util control=sequence
//ff:what buildSupersedeDoc — a Doc fixture exercising a supersede edge and level fail/review nodes
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildSupersedeDoc builds a small Doc that exercises the Phase 012 additions: a
// case with three rule nodes (one unset, one `level fail`, one `level review`)
// and a supersede edge (slotKind 4) between two of them. Blocks are laid out in
// document order with sequential non-overlapping spans that tile the Source
// exactly, so VerifyTiling holds and mode-1 byte identity is meaningful.
func buildSupersedeDoc() *doc.Doc {
	a := &ast.Document{Subject: "guard"}
	a.Sees = []ast.See{{Alias: "sec", Package: "sec/core"}}
	a.Cases = []ast.Case{{
		Name: "can enter",
		Nodes: []ast.Node{
			{Name: "has badge", Role: ast.GeneralRule, Using: &ast.Ref{Alias: "sec", Name: "hasBadge"}},
			{Name: "badge revoked", Role: ast.CounterRule, Using: &ast.Ref{Alias: "sec", Name: "isRevoked"}, Level: ast.LevelFail},
			{Name: "under review", Role: ast.ExceptRule, Using: &ast.Ref{Alias: "sec", Name: "pendingReview"}, Level: ast.LevelReview},
		},
		Supersedes: []ast.Supersede{{Superseder: "badge revoked", Superseded: "has badge"}},
	}}

	var src []byte
	var blocks []doc.Block
	add := func(kind doc.BlockKind, ref doc.StmtRef) {
		off := uint32(len(src))
		src = append(src, '-', '-', '-', '-')
		blocks = append(blocks, doc.Block{Kind: kind, Span: doc.Span{Off: off, Len: 4}, Ref: ref})
	}
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionSubject})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionSubject})
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionSee})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{0}})
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionCases})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 2}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 4, 0}})

	return &doc.Doc{AST: a, Source: src, Blocks: blocks}
}
