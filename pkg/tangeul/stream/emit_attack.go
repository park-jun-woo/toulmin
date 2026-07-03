//ff:func feature=tangl type=util control=sequence
//ff:what emitAttack — encode a `don't <target> when <attacker>` defeat edge
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitAttack encodes a defeat edge (attacker -> target, §6.2) whose endpoints are
// in-case rule nodes resolved through the pass-1 name map. An undeclared node is
// a hard error.
func emitAttack(e *enc, ref doc.StmtRef, a ast.Attack, span word.Span) ([]word.Word, error) {
	caseIdx := ref.Path[0]
	attacker, ok := e.tt.nodeByCase[nodeKey(caseIdx, a.Attacker)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: attack references undeclared node %q", a.Attacker)
	}
	target, ok := e.tt.nodeByCase[nodeKey(caseIdx, a.Target)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: attack references undeclared node %q", a.Target)
	}
	var body []word.Word
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropAttacks, Src: packet.TID(attacker), DstTID: packet.TID(target), Span: span}); err != nil {
		return nil, err
	}
	return body, nil
}
