//ff:func feature=tangl type=parser control=selection
//ff:what renderExecKo — render an execution edge (do [한 번만]/undo 되돌린다/run 진행한다)
package ko

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderExecKo renders an execution edge back to its KO surface form by kind:
// "`<node>` 면 [확신 N% <문턱>이면] `<F>` [한 번만] 실행한다" (do), "`<node>` 면 `<F>` 로 되돌린다"
// (undo), or "`<node>` 면 `<case>` 진행한다" (run). The condition node comes first; a 확신
// gate replaces the bare 면, and the undo 로(로/으로) is 받침-selected. It inverts
// parseDoExecKo/parseUndoExecKo/parseRunExecKo.
func renderExecKo(e ast.Exec) (string, error) {
	node := "`" + e.Node + "`"
	switch e.Kind {
	case ast.DoExec:
		cond, err := doCondKo(e)
		if err != nil {
			return "", err
		}
		s := node + " " + cond + " " + renderRefKo(*e.Func)
		if e.Once {
			s += " 한 번만"
		}
		return s + " 실행한다", nil
	case ast.UndoExec:
		return node + " " + josaFor(e.Node, "이면", "면") + " " + renderRefKo(*e.Func) +
			" " + josaFor(e.Func.Name, "으로", "로") + " 되돌린다", nil
	case ast.RunExec:
		return node + " " + josaFor(e.Node, "이면", "면") + " `" + e.Case + "` 진행한다", nil
	default:
		return "", fmt.Errorf("ko: unknown ExecKind %d", int(e.Kind))
	}
}
