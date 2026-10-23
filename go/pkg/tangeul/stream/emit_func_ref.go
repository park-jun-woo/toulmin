//ff:func feature=tangl type=util control=sequence
//ff:what emitFuncRef — materialize a function reference as a Func node, returning its TID
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitFuncRef materializes an ast.Ref (a `using`/`do`/`undo` function spelling)
// as a FUNC node: its F_FUNC slot binds the package Alias TID
// (0 for a local reference), VerbCode stays unbound (§6), and its name rides as a TEXT + LABEL. The name Text is appended to prefix and the node+label to body;
// the returned Func node TID is what the referring node/edge places in its own
// FuncRef Alias slot, so the reference round-trips through this node. An
// undeclared alias is a hard error.
func emitFuncRef(e *enc, prefix, body *[]core.Word, ref ast.Ref, span toulmin.Span) (toulmin.TID, error) {
	var aliasTID toulmin.TID
	if ref.Alias != "" {
		t, ok := e.tt.aliasByName[ref.Alias]
		if !ok {
			return 0, fmt.Errorf("tangeul/stream: function reference to undeclared alias %q", ref.Alias)
		}
		aliasTID = t
	}
	nameTID := addText(e, prefix, ref.Name)
	fTID := issueID(e)
	appendNode(e, body, toulmin.NodePacket{
		TID:   toulmin.TID(fTID),
		Type:  toulmin.NodeFunc,
		Flags: toulmin.NodeFlagFunc,
		Func:  &toulmin.FuncRef{Alias: toulmin.TID(aliasTID)},
		Span:  span,
	})
	appendEdge(e, body, toulmin.EdgePacket{Prop: toulmin.PropLabel, Src: fTID, DstTID: nameTID, Span: span})
	return fTID, nil
}
