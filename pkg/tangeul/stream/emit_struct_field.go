//ff:func feature=tangl type=util control=sequence
//ff:what emitStructField — encode a struct `has` field (Field node + name/type labels + has edge)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitStructField encodes a struct definition's `has f as Type` line: a Field
// node with name and type labels, and a has Edge from its owning Term node.
func emitStructField(e *enc, ref doc.StmtRef, f ast.Field, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	termTID := e.tt.byRef[tidKey(doc.StmtRef{Sec: codebook.SectionDefinitions, Path: ref.Path[:1]})]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeField, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, f.Name, span); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, f.Type, span); err != nil {
		return nil, err
	}
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropHas, Src: packet.TID(termTID), DstTID: packet.TID(self), Span: span}); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
