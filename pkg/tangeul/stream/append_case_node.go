//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what appendCaseNode — rebuild an in-case rule node (role, using, checking, with, qualifier)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// appendCaseNode rebuilds an ast.Node from a rule node packet inside the open
// case: its role, its crisp reporting level, an F_FUNC warrant resolved through
// the Func node, `with` term names from the F_WITH Text TIDs, and an F_QUAL
// qualifier. An F_CHECK case reference is deferred as a fixup (verdict
// composition may target a case declared later). It appends to the case's Nodes
// and returns the StmtRef.
func appendCaseNode(d *dec, p packet.NodePacket, name string, labels map[uint16][]string) (doc.StmtRef, error) {
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
		ref, err := resolveFuncRef(d, uint16(p.Func.Alias))
		if err != nil {
			return doc.StmtRef{}, err
		}
		n.Using = &ref
	}
	for _, wt := range p.With {
		n.With = append(n.With, d.texts[uint16(wt)])
	}
	if p.Qual != nil {
		q := p.Qual.Float()
		n.Qualified = &q
	}
	ni := len(d.ast.Cases[d.curCase].Nodes)
	d.ast.Cases[d.curCase].Nodes = append(d.ast.Cases[d.curCase].Nodes, n)
	if p.Check != nil {
		d.checkFixups = append(d.checkFixups, nameFixup{caseIdx: d.curCase, slotIdx: ni, tid: uint16(*p.Check)})
	}
	return doc.StmtRef{Sec: codebook.SectionCases, Path: []int{d.curCase, 1, ni}}, nil
}
