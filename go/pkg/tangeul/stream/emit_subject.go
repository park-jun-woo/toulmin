//ff:func feature=tangl type=util control=sequence
//ff:what emitSubject — encode the Subject declaration as a Namespace node with a name label
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// emitSubject encodes `this document is x` as a Namespace node (§5.4) whose name
// label carries the subject string verbatim.
func emitSubject(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeNamespace, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, elem.(string), span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
