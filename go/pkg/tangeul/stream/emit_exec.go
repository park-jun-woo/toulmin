//ff:func feature=tangl type=util control=selection
//ff:what emitExecStmt — encode a do/undo/run execution edge (once flag, confidence gate, dst)
package stream

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/doc"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// emitExecStmt encodes an ast.Exec as its execution Edge (§6): the PropCode
// follows the ExecKind, the src is the `when <node>` trigger, `once` and the
// confidence gate ride as flags, and the dst is a Func node (do/undo) or the
// target case TID (run). An undeclared trigger or run target is a hard error.
func emitExecStmt(e *enc, ref doc.StmtRef, x ast.Exec, span toulmin.Span) ([]core.Word, error) {
	trigger, ok := e.tt.nodeByCase[nodeKey(ref.Path[0], x.Node)]
	if !ok {
		return nil, fmt.Errorf("tangeul/stream: exec `when` references undeclared node %q", x.Node)
	}
	prop, err := codebook.ExecPropCode(x.Kind)
	if err != nil {
		return nil, err
	}
	ep, err := execFlags(x)
	if err != nil {
		return nil, err
	}
	ep.Prop, ep.Src, ep.Span = prop, toulmin.TID(trigger), span
	var pre, body []core.Word
	switch x.Kind {
	case ast.DoExec, ast.UndoExec:
		if x.Func == nil {
			return nil, fmt.Errorf("tangeul/stream: do/undo exec missing function reference")
		}
		fTID, err := emitFuncRef(e, &pre, &body, *x.Func, span)
		if err != nil {
			return nil, err
		}
		ep.DstFunc = &toulmin.FuncRef{Alias: toulmin.TID(fTID)}
	case ast.RunExec:
		dst, ok := e.tt.caseByName[x.Case]
		if !ok {
			return nil, fmt.Errorf("tangeul/stream: run references undeclared case %q", x.Case)
		}
		ep.DstTID = toulmin.TID(dst)
	}
	if err := appendEdge(e, &body, ep); err != nil {
		return nil, err
	}
	return append(pre, body...), nil
}
