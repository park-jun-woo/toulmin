//ff:func feature=tangl type=util control=sequence
//ff:what buildAlias — append a See alias, resolving its package path from the importedFrom edge
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildAlias appends the ast.See an Alias node declares: its name is the label
// and its package path is the Text the importedFrom edge points at.
func buildAlias(d *dec, name string, rel []packet.EdgePacket) (doc.StmtRef, error) {
	see := ast.See{Alias: name}
	if e, ok := findRel(rel, codebook.PropImportedFrom); ok {
		see.Package = d.texts[uint16(e.DstTID)]
	}
	idx := len(d.ast.Sees)
	d.ast.Sees = append(d.ast.Sees, see)
	return doc.StmtRef{Sec: codebook.SectionSee, Path: []int{idx}}, nil
}
