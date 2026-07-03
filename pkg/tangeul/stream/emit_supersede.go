//ff:func feature=tangl type=util control=sequence
//ff:what emitSupersede — encode a `` `superseder` supersedes `superseded` `` exclusion edge
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitSupersede encodes a supersede edge (superseder -> superseded, §6.2) whose
// endpoints are in-case rule nodes resolved through the pass-1 name map. It
// mirrors emitAttack — Src is the active side (Superseder), DstTID the excluded
// side (Superseded), with no order reversal. An undeclared node is a hard error.
func emitSupersede(e *enc, ref doc.StmtRef, s ast.Supersede, span word.Span) ([]word.Word, error) {
	caseIdx := ref.Path[0]
	superseder, ok := e.tt.nodeByCase[nodeKey(caseIdx, s.Superseder)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: supersede references undeclared node %q", s.Superseder)
	}
	superseded, ok := e.tt.nodeByCase[nodeKey(caseIdx, s.Superseded)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: supersede references undeclared node %q", s.Superseded)
	}
	var body []word.Word
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropSupersedes, Src: packet.TID(superseder), DstTID: packet.TID(superseded), Span: span}); err != nil {
		return nil, err
	}
	return body, nil
}
