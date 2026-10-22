//ff:func feature=tangl type=util control=selection
//ff:what buildExecEdge — rebuild a do/undo/run exec from its edge (once, gate, func or run target)
package stream

import (
	"errors"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// buildExecEdge rebuilds an ast.Exec from its execution edge in the open case:
// the ExecKind from the PropCode, the `when` trigger from the src, `once` and the
// confidence gate from flags, and either a do/undo Func (via its Func node) or a
// run target case name deferred as a fixup.
func buildExecEdge(d *dec, e toulmin.EdgePacket) (doc.StmtRef, error) {
	if d.curCase < 0 {
		return doc.StmtRef{}, errors.New("tangeul/stream: exec outside any case")
	}
	kind, err := codebook.ExecKindFromProp(e.Prop)
	if err != nil {
		return doc.StmtRef{}, err
	}
	x := ast.Exec{Kind: kind, Node: d.nodeName[toulmin.TID(e.Src)], Once: e.Flags&toulmin.EdgeFlagOnce != 0}
	if e.Gate != nil {
		op, err := codebook.GateOpToAST(e.Gate.Op)
		if err != nil {
			return doc.StmtRef{}, err
		}
		x.Certainty = &ast.Certainty{Op: op, Percent: int(e.Gate.Percent)}
	}
	ei := len(d.ast.Cases[d.curCase].Execs)
	switch kind {
	case ast.DoExec, ast.UndoExec:
		if e.DstFunc == nil {
			return doc.StmtRef{}, errors.New("tangeul/stream: do/undo edge missing FuncRef dst")
		}
		ref, err := resolveFuncRef(d, toulmin.TID(e.DstFunc.Alias))
		if err != nil {
			return doc.StmtRef{}, err
		}
		x.Func = &ref
	case ast.RunExec:
		d.execCaseFixups = append(d.execCaseFixups, nameFixup{caseIdx: d.curCase, slotIdx: ei, tid: toulmin.TID(e.DstTID)})
	}
	d.ast.Cases[d.curCase].Execs = append(d.ast.Cases[d.curCase].Execs, x)
	return doc.StmtRef{Sec: toulmin.SectionCases, Path: []int{d.curCase, 3, ei}}, nil
}
