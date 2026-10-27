//ff:func feature=tangl type=util control=sequence
//ff:what controlWords — issue an Edge Id and encode a SECTION / CASE-BEGIN / CASE-END / PROSE record
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// controlWords encodes a framing record with a fresh Edge Id: CASE-BEGIN /
// CASE-END carry the case Edge Id in arg, SECTION a section code, PROSE a span.
// An encoding error is recorded on e (setErr) and no words are returned.
func controlWords(e *enc, code toulmin.ControlCode, arg uint64, span *toulmin.Span) []core.Word {
	ws, err := toulmin.EncodeControl(toulmin.ControlPacket{ID: issueID(e), Code: code, Arg: arg, Span: span}, e.wire)
	setErr(e, err)
	return ws
}
