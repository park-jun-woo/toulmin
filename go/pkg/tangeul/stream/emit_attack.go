//ff:func feature=tangl type=util control=sequence
//ff:what emitAttack — encode a defeat edge as an ATTACKS LINK (F_FORCE + SIGNED-DECIMAL force when explicit)
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitAttack encodes a defeat edge (attacker -> target, §6.2) whose endpoints are
// in-case rule nodes resolved through the pass-1 name map. An explicit force
// (Attack.Force non-nil, even 1) sets F_FORCE and carries the force as a
// SIGNED-DECIMAL slot; the omitted form leaves F_FORCE off (force 1). An
// undeclared node or a force outside [0, 1] is a hard error.
func emitAttack(e *enc, ref doc.StmtRef, a ast.Attack, span toulmin.Span) ([]core.Word, error) {
	caseIdx := ref.Path[0]
	attacker, ok := e.tt.nodeByCase[nodeKey(caseIdx, a.Attacker)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: attack references undeclared node %q", a.Attacker)
	}
	target, ok := e.tt.nodeByCase[nodeKey(caseIdx, a.Target)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: attack references undeclared node %q", a.Target)
	}
	ep := toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: toulmin.TID(attacker), DstTID: toulmin.TID(target), Span: span}
	if a.Force != nil {
		v, err := toulmin.VFixFromFloat(*a.Force)
		if err != nil {
			return nil, err
		}
		ep.Flags |= toulmin.EdgeFlagForce
		ep.Force = &v
	}
	var body []core.Word
	appendEdge(e, &body, ep)
	return body, nil
}
