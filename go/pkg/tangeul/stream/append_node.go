//ff:func feature=tangl type=util control=sequence
//ff:what appendNode — encode a NODE record at the document widths and append it to the statement body
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// appendNode encodes p at the document's widths and appends the record to
// *body, the span-bearing region a statement's NODE/LINK records accumulate
// into. A node's Edge Id is its pass-1 TID, set by the caller.
func appendNode(e *enc, body *[]core.Word, p toulmin.NodePacket) error {
	w, err := toulmin.EncodeNode(p, e.wire)
	if err != nil {
		return err
	}
	*body = append(*body, w...)
	return nil
}
