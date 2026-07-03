//ff:func feature=tangl type=util control=sequence
//ff:what emitRequire — encode a `<field> is required [as Type]` statement (Field node + requires edge)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitRequire encodes an ast.Require shared by case and endpoint bodies: a Field
// node (self TID), its name label, an optional second label carrying the `as
// Type` spelling verbatim, and a requires Edge from the owner (case or endpoint).
// The owner's node type lets the reader route the field without section context.
func emitRequire(e *enc, r ast.Require, self, owner uint16, span word.Span) ([]word.Word, error) {
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeField, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, r.Field, span); err != nil {
		return nil, err
	}
	if r.Type != "" {
		if err := appendLabel(e, &pre, &body, self, r.Type, span); err != nil {
			return nil, err
		}
	}
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropRequires, Src: packet.TID(owner), DstTID: packet.TID(self), Span: span}); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
