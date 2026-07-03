//ff:func feature=tangl type=util control=sequence
//ff:what emitCaseEdge — encode an edge from an owner node to a case named target
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// emitCaseEdge encodes a single edge from owner to the case named target — the
// shared shape of endpoint runs/checks and internal runs/checks. The target
// resolves through the pass-1 case map; an undeclared case is a hard error.
func emitCaseEdge(e *enc, prop codebook.PropCode, owner uint16, target string, span word.Span) ([]word.Word, error) {
	dst, ok := e.tt.caseByName[target]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: edge references undeclared case %q", target)
	}
	var body []word.Word
	if err := appendEdge(&body, packet.EdgePacket{Prop: prop, Src: packet.TID(owner), DstTID: packet.TID(dst), Span: span}); err != nil {
		return nil, err
	}
	return body, nil
}
