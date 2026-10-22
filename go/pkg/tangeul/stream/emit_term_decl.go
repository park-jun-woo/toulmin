//ff:func feature=tangl type=util control=sequence
//ff:what emitTermDecl — encode a Definitions term (Term node, role=DefKind, means/sameAs for a const)
package stream

import (
	"encoding/json"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitTermDecl encodes a term declaration as a Term node whose role bits carry
// the DefKind, plus its name label. A const term adds a means Edge to the literal
// value and, when it has an `as` clause, a sameAs Edge to that spec reference
// materialized as a Func node. A struct term's fields are separate statements.
func emitTermDecl(e *enc, ref doc.StmtRef, d ast.Definition, span toulmin.Span) ([]core.Word, error) {
	self := e.tt.byRef[tidKey(ref)]
	var pre, body []core.Word
	if err := appendNode(e, &body, toulmin.NodePacket{TID: toulmin.TID(self), Type: toulmin.NodeTerm, Role: toulmin.RoleCode(d.Kind), Span: span}); err != nil {
		return nil, err
	}
	if err := appendLabel(e, &pre, &body, self, d.Name, span); err != nil {
		return nil, err
	}
	if d.Kind == ast.StructDef {
		return append(pre, body...), nil
	}
	value := d.Value
	if d.Kind == ast.StringListDef {
		encoded, _ := json.Marshal(d.Strings) // a []string always marshals
		value = string(encoded)
	}
	if d.Kind == ast.FileSpecDef {
		encoded, _ := json.Marshal(struct {
			LoadType string `json:"loadType"`
			FilePath string `json:"filePath"`
		}{d.LoadType, d.FilePath}) // two strings always marshal
		value = string(encoded)
	}
	valTID, err := addText(e, &pre, value)
	if err != nil {
		return nil, err
	}
	if err := appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropMeans, Src: toulmin.TID(self), DstTID: toulmin.TID(valTID), Span: span}); err != nil {
		return nil, err
	}
	if d.Kind == ast.ConstDef && d.SpecRef != nil {
		fTID, err := emitFuncRef(e, &pre, &body, *d.SpecRef, span)
		if err != nil {
			return nil, err
		}
		if err := appendEdge(e, &body, toulmin.EdgePacket{Prop: toulmin.PropSameAs, Src: toulmin.TID(self), DstTID: toulmin.TID(fTID), Span: span}); err != nil {
			return nil, err
		}
	}
	return append(pre, body...), nil
}
