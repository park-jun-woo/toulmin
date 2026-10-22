//ff:func feature=tangl type=util control=sequence
//ff:what buildTerm — append a Definitions term (DefKind from role, value and spec ref from edges)
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildTerm appends the ast.Definition a Term node declares: its DefKind from the
// role bits, a const's literal value from the means edge, and an `as` spec
// reference from the sameAs edge (resolved through its Func node).
func buildTerm(d *dec, tid toulmin.TID, name string, role toulmin.RoleCode, rel []toulmin.EdgePacket) (doc.StmtRef, error) {
	def := ast.Definition{Name: name, Kind: ast.DefKind(role)}
	if e, ok := findRel(rel, toulmin.PropMeans); ok {
		if err := applyTermValue(&def, d.texts[e.DstTID]); err != nil {
			return doc.StmtRef{}, err
		}
	}
	sp, err := termSpecRef(d, rel)
	if err != nil {
		return doc.StmtRef{}, err
	}
	def.SpecRef = sp
	idx := len(d.ast.Defs)
	d.ast.Defs = append(d.ast.Defs, def)
	d.defByTID[tid] = idx
	return doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{idx}}, nil
}
