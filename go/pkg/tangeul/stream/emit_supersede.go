//ff:func feature=tangl type=util control=sequence
//ff:what emitSupersede — encode a `` `superseder` supersedes `superseded` `` exclusion edge
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitSupersede encodes a supersede edge (superseder -> superseded, §6.2) whose
// endpoints are in-case rule nodes resolved through the pass-1 name map. It
// mirrors emitAttack — Src is the active side (Superseder), DstTID the excluded
// side (Superseded), with no order reversal. An undeclared node is a hard error.
func emitSupersede(e *enc, ref doc.StmtRef, s ast.Supersede, span toulmin.Span) ([]core.Word, error) {
	caseIdx := ref.Path[0]
	superseder, ok := e.tt.nodeByCase[nodeKey(caseIdx, s.Superseder)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: supersede references undeclared node %q", s.Superseder)
	}
	superseded, ok := e.tt.nodeByCase[nodeKey(caseIdx, s.Superseded)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: supersede references undeclared node %q", s.Superseded)
	}
	var body []core.Word
	appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropSupersedes, Src: toulmin.TID(superseder), DstTID: toulmin.TID(superseded), Span: span})
	return body, nil
}
