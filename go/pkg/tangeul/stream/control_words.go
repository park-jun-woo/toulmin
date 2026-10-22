//ff:func feature=tangl type=util control=sequence
//ff:what controlWords — issue an Edge Id and encode a SECTION / CASE-BEGIN / CASE-END / PROSE record
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
)

// controlWords encodes a framing record with a fresh Edge Id: CASE-BEGIN /
// CASE-END carry the case Edge Id in arg, SECTION a section code, PROSE a span.
func controlWords(e *enc, code toulmin.ControlCode, arg uint64, span *toulmin.Span) ([]core.Word, error) {
	return toulmin.EncodeControl(toulmin.ControlPacket{ID: issueID(e), Code: code, Arg: arg, Span: span}, e.wire)
}
