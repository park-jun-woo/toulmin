//ff:func feature=tangl type=util control=sequence
//ff:what emitInternalDecl — encode a tangl:Internal trigger head (`on`/`every … until`)
package stream

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/geul-org/gogeul/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitInternalDecl encodes the head line of a tangl:Internal entry as a Trigger
// node whose role bits carry the InternalKind. A single label holds the `on`
// event (OnEvent) or the `every` interval (EveryTick) source verbatim; an
// EveryTick with an `until` condition adds an until Edge to that case.
func emitInternalDecl(e *enc, ref doc.StmtRef, in ast.Internal, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeTrigger, Role: toulmin.RoleCode(in.Kind), Span: span})
	head := in.Event
	if in.Kind == ast.EveryTick {
		head = in.Interval
	}
	appendLabel(e, &pre, &body, self, head, span)
	if in.Kind == ast.EveryTick && in.Until != "" {
		dst, ok := e.tt.caseByName[in.Until]
		if !ok {
			return nil, fmt.Errorf("tangeul/stream: internal until references undeclared case %q", in.Until)
		}
		appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropUntil, Src: toulmin.TID(self), DstTID: toulmin.TID(dst), Span: span})
	}
	return append(pre, body...), nil
}
