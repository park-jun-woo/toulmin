//ff:func feature=tangl type=util control=sequence
//ff:what emitInternalDecl — encode a tangl:Internal trigger head (`on`/`every … until`)
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitInternalDecl encodes the head line of a tangl:Internal entry as a Trigger
// node whose role bits carry the InternalKind. A single label holds the `on`
// event (OnEvent) or the `every` interval (EveryTick) source verbatim; an
// EveryTick with an `until` condition adds an until Edge to that case.
func emitInternalDecl(e *enc, ref doc.StmtRef, in ast.Internal, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeTrigger, Role: codebook.RoleCode(in.Kind), Span: span}); err != nil {
		return nil, err
	}
	head := in.Event
	if in.Kind == ast.EveryTick {
		head = in.Interval
	}
	if err := appendLabel(e, &pre, &body, self, head, span); err != nil {
		return nil, err
	}
	if in.Kind == ast.EveryTick && in.Until != "" {
		dst, ok := e.tt.caseByName[in.Until]
		if !ok {
			return nil, fmt.Errorf("tangeul/stream: internal until references undeclared case %q", in.Until)
		}
		if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropUntil, Src: packet.TID(self), DstTID: packet.TID(dst), Span: span}); err != nil {
			return nil, err
		}
	}
	return append(pre, body...), nil
}
