//ff:func feature=tangl type=util control=sequence
//ff:what emitCaseDecl — encode an `in case of x` declaration as a Case node with a name label
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitCaseDecl encodes a case declaration as a Case node and its name label. The
// enclosing CASE_BEGIN/END controls (emitted by WriteDoc) bound the case body.
func emitCaseDecl(e *enc, ref doc.StmtRef, c ast.Case, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeCase, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, c.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
