//ff:func feature=tangl type=util control=sequence
//ff:what emitEndpointDecl — encode a `provides x` endpoint as an Endpoint node with a name label
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitEndpointDecl encodes an endpoint declaration as an Endpoint node and its
// name label. Its requires/runs/checks are separate statements routed by this
// endpoint's TID as the edge src.
func emitEndpointDecl(e *enc, ref doc.StmtRef, ep ast.Endpoint, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeEndpoint, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, ep.Name, span); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
