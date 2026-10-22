//ff:func feature=tangl type=util control=sequence
//ff:what emitStructField — encode a struct `has` field (Field node + name/type labels + has edge)
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitStructField encodes a struct definition's `has f as Type` line: a Field
// node with name and type labels, and a has Edge from its owning Term node.
func emitStructField(e *enc, ref doc.StmtRef, f ast.Field, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	termTID := e.tt.byRef[tidKey(doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: ref.Path[:1]})]
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeField, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, f.Name, span); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, f.Type, span); err != nil {
		return nil, err
	}
	if err := appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropHas, Src: toulmin.TID(termTID), DstTID: toulmin.TID(self), Span: span}); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
