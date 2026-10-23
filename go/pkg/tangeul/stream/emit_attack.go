//ff:func feature=tangl type=util control=sequence
//ff:what emitAttack — encode a `don't <target> when <attacker>` defeat edge
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitAttack encodes a defeat edge (attacker -> target, §6.2) whose endpoints are
// in-case rule nodes resolved through the pass-1 name map. An undeclared node is
// a hard error.
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
	var body []core.Word
	appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropAttacks, Src: toulmin.TID(attacker), DstTID: toulmin.TID(target), Span: span})
	return body, nil
}
