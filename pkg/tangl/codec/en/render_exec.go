//ff:func feature=tangl type=parser control=selection
//ff:what renderExec — render an execution edge (do [once]/undo/run ... when `node` [certainty])
package en

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderExec renders an execution edge back to its surface form by kind:
// "do <ref> [once] when `<node>` [if <op> N% certain]", "undo <ref> when `<node>`",
// or "run `<case>` when `<node>`". It inverts parseDoExec/parseUndoExec/parseRunExec.
func renderExec(e ast.Exec) (string, error) {
	switch e.Kind {
	case ast.DoExec:
		s := "do " + renderRef(*e.Func)
		if e.Once {
			s += " once"
		}
		s += " when `" + e.Node + "`"
		if e.Certainty != nil {
			s += " " + renderCertainty(*e.Certainty)
		}
		return s, nil
	case ast.UndoExec:
		return "undo " + renderRef(*e.Func) + " when `" + e.Node + "`", nil
	case ast.RunExec:
		return "run `" + e.Case + "` when `" + e.Node + "`", nil
	default:
		return "", fmt.Errorf("en: unknown ExecKind %d", int(e.Kind))
	}
}
