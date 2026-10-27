//ff:func feature=tangl type=util control=sequence
//ff:what appendLabel — add a TEXT for s to the prefix and a LABEL record (src → that TEXT) to the body
package stream

import (
	"github.com/geul-org/gogeul/pkg/core"
	"github.com/geul-org/gogeul/pkg/toulmin"
)

// appendLabel adds a TEXT record for s to the statement prefix and a LABEL
// record (src → that TEXT) to the body — the standard way a node's spelling is
// bound losslessly. The TEXT rides in the prefix so grouping stays
// prefix-then-body.
func appendLabel(e *enc, pre, body *[]core.Word, src toulmin.TID, s string, span toulmin.Span) {
	appendEdge(e, body, toulmin.EdgePacket{Prop: toulmin.PropLabel, Src: src, DstTID: addText(e, pre, s), Span: span})
}
