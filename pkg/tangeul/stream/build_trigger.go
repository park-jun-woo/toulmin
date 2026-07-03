//ff:func feature=tangl type=util control=sequence
//ff:what buildTrigger — append a tangl:Internal trigger (kind from role, event/interval, until edge)
package stream

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildTrigger appends the ast.Internal a Trigger node declares: its kind from
// the role bits, the head label as the `on` event or `every` interval, and an
// `until` case from the until edge. Runs/checks are separate statements routed by
// this trigger's TID.
func buildTrigger(d *dec, tid uint16, name string, role codebook.RoleCode, rel []packet.EdgePacket) (doc.StmtRef, error) {
	in := ast.Internal{Kind: ast.InternalKind(role)}
	if in.Kind == ast.OnEvent {
		in.Event = name
	} else {
		in.Interval = name
	}
	if e, ok := findRel(rel, codebook.PropUntil); ok {
		in.Until = d.nodeName[uint16(e.DstTID)]
	}
	idx := len(d.ast.Internals)
	d.ast.Internals = append(d.ast.Internals, in)
	d.trByTID[tid] = idx
	return doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{idx}}, nil
}
