//ff:func feature=tangl type=util control=sequence
//ff:what buildField — append a Field the node declares, routing by its has or requires edge
package stream

import (
	"errors"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildField appends the Field a NodeField declares: a has edge makes it a struct
// definition field (owner = its term), a requires edge makes it a required field
// routed to its case or endpoint owner. The optional second label is the type.
func buildField(d *dec, tid toulmin.TID, name string, labels map[toulmin.TID][]string, rel []toulmin.EdgePacket) (doc.StmtRef, error) {
	typ := ""
	if ls := labels[tid]; len(ls) > 1 {
		typ = ls[1]
	}
	if e, ok := findRel(rel, toulmin.PropHas); ok {
		di := d.defByTID[toulmin.TID(e.Src)]
		fi := len(d.ast.Defs[di].Fields)
		d.ast.Defs[di].Fields = append(d.ast.Defs[di].Fields, ast.Field{Name: name, Type: typ})
		return doc.StmtRef{Sec: toulmin.SectionDefinitions, Path: []int{di, fi}}, nil
	}
	if e, ok := findRel(rel, toulmin.PropRequires); ok {
		return appendRequire(d, toulmin.TID(e.Src), ast.Require{Field: name, Type: typ})
	}
	return doc.StmtRef{}, errors.New("tangeul/stream: field node has no has/requires link")
}
