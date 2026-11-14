//ff:func feature=tangl type=util control=sequence
//ff:what buildAlias — append a See alias, resolving its package path from the importedFrom edge
package stream

import (
	"github.com/geul-org/gogeul/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildAlias appends the ast.See an Alias node declares: its name is the label
// and its package path is the Text the importedFrom edge points at.
func buildAlias(d *dec, name string, rel []toulmin.EdgePacket) (doc.StmtRef, error) {
	see := ast.See{Alias: name}
	if e, ok := findRel(rel, toulmin.PropImportedFrom); ok {
		see.Package = d.texts[toulmin.TID(e.DstTID)]
	}
	idx := len(d.ast.Sees)
	d.ast.Sees = append(d.ast.Sees, see)
	return doc.StmtRef{Sec: toulmin.SectionSee, Path: []int{idx}}, nil
}
