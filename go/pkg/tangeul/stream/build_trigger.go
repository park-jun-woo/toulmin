//ff:func feature=tangl type=util control=sequence
//ff:what buildTrigger — append a tangl:Internal trigger (kind from role, event/interval, until edge)
package stream

import (
	"github.com/geul-org/geul/go/pkg/toulmin"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildTrigger appends the ast.Internal a Trigger node declares: its kind from
// the role bits, the head label as the `on` event or `every` interval, and an
// `until` case from the until edge. Runs/checks are separate statements routed by
// this trigger's TID.
func buildTrigger(d *dec, tid toulmin.TID, name string, role toulmin.RoleCode, rel []toulmin.EdgePacket) (doc.StmtRef, error) {
	in := ast.Internal{Kind: ast.InternalKind(role)}
	if in.Kind == ast.OnEvent {
		in.Event = name
	} else {
		in.Interval = name
	}
	if e, ok := findRel(rel, toulmin.PropUntil); ok {
		in.Until = d.nodeName[toulmin.TID(e.DstTID)]
	}
	idx := len(d.ast.Internals)
	d.ast.Internals = append(d.ast.Internals, in)
	d.trByTID[tid] = idx
	return doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{idx}}, nil
}
