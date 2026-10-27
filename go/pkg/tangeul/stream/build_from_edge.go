//ff:func feature=tangl type=util control=selection
//ff:what buildFromEdge — append the AST element an edge-only statement declares and return its StmtRef
package stream

import (
	"errors"
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildFromEdge handles a statement whose only relation is a single edge (attack,
// supersede, exec, endpoint run/check, internal run/check) and returns its
// StmtRef. Attacks, supersedes and execs belong to the open case; runs and checks
// route through their owner's TID. Exec building (with its forward `run` fixup) is
// delegated to buildExecEdge.
func buildFromEdge(d *dec, rel []toulmin.EdgePacket) (doc.StmtRef, error) {
	if len(rel) != 1 {
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: edge statement must have exactly one relation edge, got %d", len(rel))
	}
	e := rel[0]
	switch e.Prop {
	case toulmin.PropAttacks:
		if d.curCase < 0 {
			return doc.StmtRef{}, errors.New("tangeul/stream: attack outside any case")
		}
		a := ast.Attack{Target: d.nodeName[toulmin.TID(e.DstTID)], Attacker: d.nodeName[toulmin.TID(e.Src)]}
		ai := len(d.ast.Cases[d.curCase].Attacks)
		d.ast.Cases[d.curCase].Attacks = append(d.ast.Cases[d.curCase].Attacks, a)
		return doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{d.curCase, 2, ai}}, nil
	case toulmin.PropSupersedes:
		if d.curCase < 0 {
			return doc.StmtRef{}, errors.New("tangeul/stream: supersede outside any case")
		}
		s := ast.Supersede{Superseder: d.nodeName[toulmin.TID(e.Src)], Superseded: d.nodeName[toulmin.TID(e.DstTID)]}
		si := len(d.ast.Cases[d.curCase].Supersedes)
		d.ast.Cases[d.curCase].Supersedes = append(d.ast.Cases[d.curCase].Supersedes, s)
		return doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{d.curCase, 4, si}}, nil
	case toulmin.PropDo, toulmin.PropUndo, toulmin.PropRun:
		return buildExecEdge(d, e)
	case toulmin.PropRuns:
		ei := d.epByTID[toulmin.TID(e.Src)]
		ri := len(d.ast.Provides[ei].Runs)
		d.ast.Provides[ei].Runs = append(d.ast.Provides[ei].Runs, d.nodeName[toulmin.TID(e.DstTID)])
		return doc.StmtRef{Sec: toulmin.SectionProvides, Path: []int{ei, 1, ri}}, nil
	case toulmin.PropEvaluates:
		return buildCheckEdge(d, e)
	case toulmin.PropTriggers:
		ii := d.trByTID[toulmin.TID(e.Src)]
		ri := len(d.ast.Internals[ii].Runs)
		d.ast.Internals[ii].Runs = append(d.ast.Internals[ii].Runs, d.nodeName[toulmin.TID(e.DstTID)])
		return doc.StmtRef{Sec: toulmin.SectionInternal, Path: []int{ii, 1, ri}}, nil
	default:
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: unexpected relation edge 0x%04X", uint16(e.Prop))
	}
}
