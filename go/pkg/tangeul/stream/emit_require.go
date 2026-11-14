//ff:func feature=tangl type=util control=sequence
//ff:what emitRequire — encode a `<field> is required [as Type]` statement (Field node + requires edge)
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitRequire encodes an ast.Require shared by case and endpoint bodies: a Field
// node (self TID), its name label, an optional second label carrying the `as
// Type` spelling verbatim, and a requires Edge from the owner (case or endpoint).
// The owner's node type lets the reader route the field without section context.
func emitRequire(e *enc, r ast.Require, self, owner toulmin.TID, span toulmin.Span) ([]core.Word, error) {
	var pre, body []core.Word
	appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeField, Span: span})
	appendLabel(e, &pre, &body, self, r.Field, span)
	if r.Type != "" {
		appendLabel(e, &pre, &body, self, r.Type, span)
	}
	appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropRequires, Src: toulmin.TID(owner), DstTID: toulmin.TID(self), Span: span})
	return append(pre, body...), nil
}
