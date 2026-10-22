//ff:func feature=cli type=command control=sequence
//ff:what buildLintErrorDoc — a Doc that both fails validation (undo without do) and lints (unguarded do reachable from Internal)
package tanglcli

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildLintErrorDoc builds a Doc that exercises runCheck's failure branch: node
// n2 carries an `undo` with no preceding `do` on that node (a validation error),
// while node n1's `do` has no `once` guard and the case is reachable from
// tangl:Internal (a lint warning). Both survive Save/Load — the stream writer
// resolves node/func refs but leaves do/undo pairing and lint to validate, which
// runs over the recovered AST.
func buildLintErrorDoc() *doc.Doc {
	a := &ast.Document{Subject: "bad"}
	a.Cases = []ast.Case{{
		Name:  "c1",
		Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}, {Name: "n2", Role: ast.GeneralRule}},
		Execs: []ast.Exec{
			{Kind: ast.DoExec, Func: &ast.Ref{Name: "report"}, Node: "n1"},
			{Kind: ast.UndoExec, Func: &ast.Ref{Name: "rollback"}, Node: "n2"},
		},
	}}
	a.Internals = []ast.Internal{{Kind: ast.OnEvent, Event: "x.y", Runs: []string{"c1"}}}

	var src []byte
	var blocks []doc.Block
	add := func(kind doc.BlockKind, ref doc.StmtRef) {
		off := uint32(len(src))
		src = append(src, '-', '-', '-', '-')
		blocks = append(blocks, doc.Block{Kind: kind, Span: doc.Span{Off: off, Len: 4}, Ref: ref})
	}
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionSubject})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionSubject})
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionCases})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 1, 1}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 3, 0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{0, 3, 1}})
	add(doc.SectionBlock, doc.StmtRef{Sec: toulmin.SectionInternal})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0}})
	add(doc.StmtBlock, doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{0, 1, 0}})

	return &doc.Doc{AST: a, Source: src, Blocks: blocks, Locale: "en"}
}
