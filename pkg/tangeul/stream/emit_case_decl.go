//ff:func feature=tangl type=util control=sequence
//ff:what emitCaseDecl — encode an `in case of x` declaration as a Case node with a name label
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitCaseDecl encodes a case declaration as a Case node and its name label. The
// enclosing CASE_BEGIN/END controls (emitted by WriteDoc) bound the case body.
func emitCaseDecl(e *enc, ref doc.StmtRef, c ast.Case, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeCase, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, c.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
