//ff:func feature=tangl type=codegen control=sequence
//ff:what writeUndoPush — writes one PushCompensation statement for an undo edge
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// writeUndoPush writes a single tangl.PushCompensation statement for one
// "undo" edge, armed at the point it is called — the caller is
// responsible for calling this immediately after the do(s) it compensates
// have already succeeded, per the spec's undo-arming rule.
func writeUndoPush(w *strings.Builder, e ast.Exec) {
	fmt.Fprintln(w, "\t\ttangl.PushCompensation(t.Ctx(), func() error {")
	fmt.Fprintf(w, "\t\t\treturn %s(t.Ctx())\n", refSelector(e.Func))
	fmt.Fprintln(w, "\t\t})")
}
