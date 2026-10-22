//ff:func feature=tangl type=util control=sequence
//ff:what appendCaseNode — rebuild an in-case rule node (role, using, checking, with, qualifier)
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// appendCaseNode rebuilds an ast.Node from a rule node packet inside the open
// case: its role, its crisp reporting level, an F_FUNC warrant resolved through
// the Func node, `with` term names from the F_WITH Text TIDs, and an F_QUAL
// qualifier. An F_CHECK case reference is deferred as a fixup (verdict
// composition may target a case declared later). It appends to the case's Nodes
// and returns the StmtRef.
func appendCaseNode(d *dec, p toulmin.NodePacket, name string, labels map[toulmin.TID][]string) (doc.StmtRef, error) {
	role, err := codebook.RoleToAST(p.Role)
	if err != nil {
		return doc.StmtRef{}, err
	}
	level, err := codebook.LevelToAST(p.Level)
	if err != nil {
		return doc.StmtRef{}, err
	}
	n := ast.Node{Name: name, Role: role, Level: level}
	if p.Func != nil {
		ref, err := resolveFuncRef(d, toulmin.TID(p.Func.Alias))
		if err != nil {
			return doc.StmtRef{}, err
		}
		n.Using = &ref
	}
	n.With, n.Grounds = textsOf(d, p.With), textsOf(d, p.Grounds)
	if p.Qual != nil {
		q := p.Qual.Float()
		n.Qualified = &q
	}
	ni := len(d.ast.Cases[d.curCase].Nodes)
	d.ast.Cases[d.curCase].Nodes = append(d.ast.Cases[d.curCase].Nodes, n)
	if p.Check != nil {
		d.checkFixups = append(d.checkFixups, nameFixup{caseIdx: d.curCase, slotIdx: ni, tid: toulmin.TID(*p.Check)})
	}
	return doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{d.curCase, 1, ni}}, nil
}
