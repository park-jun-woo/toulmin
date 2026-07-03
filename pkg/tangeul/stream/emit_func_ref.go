//ff:func feature=tangl type=util control=sequence
//ff:what emitFuncRef — materialize a function reference as a Func node, returning its TID
package stream

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitFuncRef materializes an ast.Ref (a `using`/`do`/`undo` function spelling)
// as a Func node (NodeFunc, §2.1): its F_FUNC slot binds the package Alias TID
// (0 for a local reference), VerbCode stays unbound (§6), and its name rides as
// a label Text. The name Text is appended to prefix and the node+label to body;
// the returned Func node TID is what the referring node/edge places in its own
// FuncRef Alias slot, so the reference round-trips through this node. An
// undeclared alias is a hard error.
func emitFuncRef(e *enc, prefix, body *[]word.Word, ref ast.Ref, span word.Span) (uint16, error) {
	var aliasTID uint16
	if ref.Alias != "" {
		t, ok := e.tt.aliasByName[ref.Alias]
		if !ok {
			return 0, fmt.Errorf("tangeul/stream: function reference to undeclared alias %q", ref.Alias)
		}
		aliasTID = t
	}
	nameTID := addText(e, prefix, ref.Name)
	fTID := e.nodeNext
	e.nodeNext++
	node, err := packet.EncodeNode(packet.NodePacket{
		TID:   packet.TID(fTID),
		Type:  codebook.NodeFunc,
		Flags: codebook.NodeFlagFunc,
		Func:  &packet.FuncRef{Alias: packet.TID(aliasTID)},
		Span:  span,
	})
	if err != nil {
		return 0, err
	}
	*body = append(*body, node...)
	lbl, err := labelEdge(fTID, nameTID, span)
	if err != nil {
		return 0, err
	}
	*body = append(*body, lbl...)
	return fTID, nil
}
