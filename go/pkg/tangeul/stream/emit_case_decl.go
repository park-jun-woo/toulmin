//ff:func feature=tangl type=util control=sequence
//ff:what emitCaseDecl — encode an `in case of x` declaration as a Case node with a name label
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitCaseDecl encodes a case declaration as a Case node and its name label. The
// enclosing CASE_BEGIN/END controls (emitted by WriteDoc) bound the case body.
func emitCaseDecl(e *enc, ref doc.StmtRef, c ast.Case, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeCase, Span: span})
	appendLabel(e, &pre, &body, self, c.Name, span)
	return append(pre, body...), nil
}
