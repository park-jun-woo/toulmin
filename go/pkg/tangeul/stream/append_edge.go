//ff:func feature=tangl type=util control=sequence
//ff:what appendEdge — issue an Edge Id for a LINK record, encode it and append it to the statement body
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// appendEdge issues the LINK record's Edge Id (relations are not referenced, so
// their ids come after pass 1), encodes ep at the document's widths and appends
// it to *body.
func appendEdge(e *enc, body *[]core.Word, ep toulmin.EdgePacket) error {
	ep.ID = issueID(e)
	w, err := toulmin.EncodeEdge(ep, e.wire)
	if err != nil {
		return err
	}
	*body = append(*body, w...)
	return nil
}
