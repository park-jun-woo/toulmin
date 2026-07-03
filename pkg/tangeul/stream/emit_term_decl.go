//ff:func feature=tangl type=util control=sequence
//ff:what emitTermDecl — encode a Definitions term (Term node, role=DefKind, means/sameAs for a const)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitTermDecl encodes a term declaration as a Term node whose role bits carry
// the DefKind, plus its name label. A const term adds a means Edge to the literal
// value and, when it has an `as` clause, a sameAs Edge to that spec reference
// materialized as a Func node. A struct term's fields are separate statements.
func emitTermDecl(e *enc, ref doc.StmtRef, d ast.Definition, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeTerm, Role: codebook.RoleCode(d.Kind), Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, d.Name, span); err != nil {
		return nil, err
	}
	if d.Kind != ast.ConstDef {
		return append(pre, body...), nil
	}
	valTID := addText(e, &pre, d.Value)
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropMeans, Src: packet.TID(self), DstTID: packet.TID(valTID), Span: span}); err != nil {
		return nil, err
	}
	if d.SpecRef != nil {
		fTID, err := emitFuncRef(e, &pre, &body, *d.SpecRef, span)
		if err != nil {
			return nil, err
		}
		if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropSameAs, Src: packet.TID(self), DstTID: packet.TID(fTID), Span: span}); err != nil {
			return nil, err
		}
	}
	return append(pre, body...), nil
}
