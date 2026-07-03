//ff:func feature=tangl type=util control=sequence
//ff:what emitEndpointDecl — encode a `provides x` endpoint as an Endpoint node with a name label
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitEndpointDecl encodes an endpoint declaration as an Endpoint node and its
// name label. Its requires/runs/checks are separate statements routed by this
// endpoint's TID as the edge src.
func emitEndpointDecl(e *enc, ref doc.StmtRef, ep ast.Endpoint, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeEndpoint, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, ep.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
