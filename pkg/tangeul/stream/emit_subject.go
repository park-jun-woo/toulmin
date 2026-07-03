//ff:func feature=tangl type=util control=sequence
//ff:what emitSubject — encode the Subject declaration as a Namespace node with a name label
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitSubject encodes `this document is x` as a Namespace node (§5.4) whose name
// label carries the subject string verbatim.
func emitSubject(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeNamespace, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, elem.(string), span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
