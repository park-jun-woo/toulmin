//ff:func feature=tangl type=util control=sequence
//ff:what emitSee — encode a `see a from path` as an Alias node with name label and importedFrom edge
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// emitSee encodes a package reference as an Alias node (§5.4): its name label is
// the alias, and an importedFrom Edge binds it to the package path carried as a
// Text payload.
func emitSee(e *enc, ref doc.StmtRef, elem any, span word.Span) ([]word.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	see := elem.(ast.See)
	var pre, body []word.Word
	if err := appendNode(&body, packet.NodePacket{TID: packet.TID(self), Type: codebook.NodeAlias, Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, see.Alias, span); err != nil {
		return nil, err
	}
	pathTID := addText(e, &pre, see.Package)
	if err := appendEdge(&body, packet.EdgePacket{Prop: codebook.PropImportedFrom, Src: packet.TID(self), DstTID: packet.TID(pathTID), Span: span}); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
