//ff:func feature=tangl type=util control=sequence
//ff:what buildAttack — rebuild an ast.Attack from an ATTACKS LINK (F_FORCE → explicit Force)
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildAttack rebuilds a defeat edge from its ATTACKS LINK: Src is the
// attacker, DstTID the target, and an F_FORCE slot (Force non-nil) becomes the
// explicit Attack.Force; without F_FORCE the edge is the omitted form (force 1).
func buildAttack(d *dec, e toulmin.EdgePacket) ast.Attack {
	a := ast.Attack{Target: d.nodeName[toulmin.TID(e.DstTID)], Attacker: d.nodeName[toulmin.TID(e.Src)]}
	if e.Force != nil {
		f := e.Force.Float()
		a.Force = &f
	}
	return a
}
