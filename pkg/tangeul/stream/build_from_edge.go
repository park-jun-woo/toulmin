//ff:func feature=tangl type=util control=selection
//ff:what buildFromEdge — append the AST element an edge-only statement declares and return its StmtRef
package stream

import (
	"errors"
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// buildFromEdge handles a statement whose only relation is a single edge (attack,
// exec, endpoint run/check, internal run/check) and returns its StmtRef. Attacks
// and execs belong to the open case; runs and checks route through their owner's
// TID. Exec building (with its forward `run` fixup) is delegated to buildExecEdge.
func buildFromEdge(d *dec, rel []packet.EdgePacket) (doc.StmtRef, error) {
	if len(rel) != 1 {
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: edge statement must have exactly one relation edge, got %d", len(rel))
	}
	e := rel[0]
	switch e.Prop {
	case codebook.PropAttacks:
		if d.curCase < 0 {
			return doc.StmtRef{}, errors.New("tangeul/stream: attack outside any case")
		}
		a := ast.Attack{Target: d.nodeName[uint16(e.DstTID)], Attacker: d.nodeName[uint16(e.Src)]}
		ai := len(d.ast.Cases[d.curCase].Attacks)
		d.ast.Cases[d.curCase].Attacks = append(d.ast.Cases[d.curCase].Attacks, a)
		return doc.StmtRef{Sec: codebook.SectionCases, Path: []int{d.curCase, 2, ai}}, nil
	case codebook.PropDo, codebook.PropUndo, codebook.PropRun:
		return buildExecEdge(d, e)
	case codebook.PropRuns:
		ei := d.epByTID[uint16(e.Src)]
		ri := len(d.ast.Provides[ei].Runs)
		d.ast.Provides[ei].Runs = append(d.ast.Provides[ei].Runs, d.nodeName[uint16(e.DstTID)])
		return doc.StmtRef{Sec: codebook.SectionProvides, Path: []int{ei, 1, ri}}, nil
	case codebook.PropEvaluates:
		return buildCheckEdge(d, e)
	case codebook.PropTriggers:
		ii := d.trByTID[uint16(e.Src)]
		ri := len(d.ast.Internals[ii].Runs)
		d.ast.Internals[ii].Runs = append(d.ast.Internals[ii].Runs, d.nodeName[uint16(e.DstTID)])
		return doc.StmtRef{Sec: codebook.SectionInternal, Path: []int{ii, 1, ri}}, nil
	default:
		return doc.StmtRef{}, fmt.Errorf("tangeul/stream: unexpected relation edge 0x%04X", uint16(e.Prop))
	}
}
