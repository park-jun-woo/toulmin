//ff:func feature=tangl type=util control=sequence
//ff:what addText — append a TEXT record carrying s to the statement prefix and return its Edge Id
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// addText issues the next Edge Id, appends a TEXT record carrying s to the
// statement prefix, and returns the Edge Id that label/means/with references
// point at. The id is written in the record, so the reader resolves it directly.
func addText(e *enc, prefix *[]core.Word, s string) (toulmin.TID, error) {
	id := issueID(e)
	ws, err := toulmin.EncodeText(toulmin.TextPacket{TID: id, Bytes: []byte(s)}, e.wire)
	if err != nil {
		return 0, err
	}
	*prefix = append(*prefix, ws...)
	return id, nil
}
