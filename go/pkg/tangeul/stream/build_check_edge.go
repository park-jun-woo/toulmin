//ff:func feature=tangl type=util control=sequence
//ff:what buildCheckEdge — route an evaluates edge to an endpoint check or an internal check
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
)

// buildCheckEdge routes an evaluates edge (PropEvaluates) by its src node type: an
// endpoint check appends to Provides[i].Checks, an internal (trigger) check
// appends to Internals[i].Checks. The dst is the checked case name.
func buildCheckEdge(d *dec, e toulmin.EdgePacket) (doc.StmtRef, error) {
	src := toulmin.TID(e.Src)
	if d.nodeType[src] == toulmin.NodeEndpoint {
		ei := d.epByTID[src]
		ci := len(d.ast.Provides[ei].Checks)
		d.ast.Provides[ei].Checks = append(d.ast.Provides[ei].Checks, d.nodeName[toulmin.TID(e.DstTID)])
		return doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei, 2, ci}}, nil
	}
	ii := d.trByTID[src]
	ci := len(d.ast.Internals[ii].Checks)
	d.ast.Internals[ii].Checks = append(d.ast.Internals[ii].Checks, d.nodeName[toulmin.TID(e.DstTID)])
	return doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{ii, 2, ci}}, nil
}
