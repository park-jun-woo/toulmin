//ff:func feature=tangl type=util control=sequence
//ff:what emitSubject — encode the Subject declaration as a Namespace node with a name label
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// emitSubject encodes `this document is x` as a Namespace node (§5.4) whose name
// label carries the subject string verbatim.
func emitSubject(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeNamespace, Span: span})
	appendLabel(e, &pre, &body, self, elem.(string), span)
	return append(pre, body...), nil
}
