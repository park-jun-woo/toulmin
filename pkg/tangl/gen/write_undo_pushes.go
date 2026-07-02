//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeUndoPushes — writes one PushCompensation statement per undo edge on a node
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeUndoPushes writes one tangl.PushCompensation statement per "undo"
// edge on this node, in document order, armed unconditionally once
// control reaches this point in the handler — i.e. after every "do" on
// the node has already succeeded, per the spec's undo-arming rule
// ("arm immediately after success").
func writeUndoPushes(w *strings.Builder, execs []ast.Exec) {
	for _, e := range execs {
		if e.Kind != ast.UndoExec {
			continue
		}
		fmt.Fprintln(w, "\t\ttangl.PushCompensation(t.Ctx(), func() error {")
		fmt.Fprintf(w, "\t\t\treturn %s(t.Ctx())\n", refSelector(e.Func))
		fmt.Fprintln(w, "\t\t})")
	}
}
