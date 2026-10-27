//ff:func feature=tangl type=util control=sequence
//ff:what appendNode — encode a NODE record at the document SPAN width and append it to the statement body
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// appendNode encodes p at the document's SPAN width and appends the record to
// *body, the span-bearing region a statement's NODE/LINK records accumulate
// into. A node's Edge Id is its pass-1 TID, set by the caller. An encoding
// error is recorded on e (setErr); the failed encoder returns no words, so
// nothing is appended.
func appendNode(e *enc, body *[]core.Word, p toulmin.NodePacket) {
	w, err := toulmin.EncodeNode(p, e.wire)
	setErr(e, err)
	*body = append(*body, w...)
}
