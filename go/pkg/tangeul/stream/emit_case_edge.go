//ff:func feature=tangl type=util control=sequence
//ff:what emitCaseEdge — encode an edge from an owner node to a case named target
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
)

// emitCaseEdge encodes a single edge from owner to the case named target — the
// shared shape of endpoint runs/checks and internal runs/checks. The target
// resolves through the pass-1 case map; an undeclared case is a hard error.
func emitCaseEdge(e *enc, prop toulmin.PropCode, owner toulmin.TID, target string, span toulmin.Span) ([]core.Word, error) {
	dst, ok := e.tt.caseByName[target]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: edge references undeclared case %q", target)
	}
	var body []core.Word
	appendEdge(e, &body, toulmin.EdgePacket{Prop: prop, Src: toulmin.TID(owner), DstTID: toulmin.TID(dst), Span: span})
	return body, nil
}
