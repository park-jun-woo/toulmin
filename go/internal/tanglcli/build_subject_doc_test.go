//ff:func feature=cli type=command control=sequence
//ff:what buildSubjectDoc — a minimal valid Doc declaring the given subject (for subject-scan tests)
package tanglcli

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildSubjectDoc returns a minimal Doc whose Subject is the given name and which
// has one Case with one node, laid out with tiling spans so Save/Load succeed. It
// is the fixture for findDocumentBySubject's directory-scan tests, replacing the
// old .md text builder now that the scanner walks .geul files.
func buildSubjectDoc(subject string) *doc.Doc {
	a := &ast.Document{Subject: subject}
	a.Cases = []ast.Case{{Name: "c1", Nodes: []ast.Node{{Name: "n1", Role: ast.GeneralRule}}}}

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

	return &doc.Doc{AST: a, Source: src, Blocks: blocks, Locale: "en"}
}
