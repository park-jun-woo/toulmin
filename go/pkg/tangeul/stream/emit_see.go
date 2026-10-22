//ff:func feature=tangl type=util control=sequence
//ff:what emitSee — encode a `see a from path` as an Alias node with name label and importedFrom edge
package stream

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitSee encodes a package reference as an Alias node (§5.4): its name label is
// the alias, and an importedFrom Edge binds it to the package path carried as a
// Text payload.
func emitSee(e *enc, ref doc.StmtRef, elem any, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	see := elem.(ast.See)
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeAlias, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, see.Alias, span); err != nil {
		return nil, err
	}
	pathTID, err := addText(e, &pre, see.Package)
	if err != nil {
		return nil, err
	}
	if err := appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropImportedFrom, Src: toulmin.TID(self), DstTID: toulmin.TID(pathTID), Span: span}); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
