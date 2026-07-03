//ff:func feature=tangl type=util control=sequence
//ff:what buildCheckEdge — route an evaluates edge to an endpoint check or an internal check
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
)

// buildCheckEdge routes an evaluates edge (PropEvaluates) by its src node type: an
// endpoint check appends to Provides[i].Checks, an internal (trigger) check
// appends to Internals[i].Checks. The dst is the checked case name.
func buildCheckEdge(d *dec, e packet.EdgePacket) (doc.StmtRef, error) {
	src := uint16(e.Src)
	if d.nodeType[src] == codebook.NodeEndpoint {
		ei := d.epByTID[src]
		ci := len(d.ast.Provides[ei].Checks)
		d.ast.Provides[ei].Checks = append(d.ast.Provides[ei].Checks, d.nodeName[uint16(e.DstTID)])
		return doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{ei, 2, ci}}, nil
	}
	ii := d.trByTID[src]
	ci := len(d.ast.Internals[ii].Checks)
	d.ast.Internals[ii].Checks = append(d.ast.Internals[ii].Checks, d.nodeName[uint16(e.DstTID)])
	return doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{ii, 2, ci}}, nil
}
