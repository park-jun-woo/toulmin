//ff:func feature=tangl type=codegen control=sequence
//ff:what renderDoStatement — writes one gated do-edge action call
package gen

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderDoStatement writes one "do" edge's gated action call: optionally
// wrapped in a certainty threshold check (self.Verdict), optionally
// wrapped in a tangl.OnceDone/OnceMark guard consumed only after the
// action succeeds, calling the leaf action as fn(t.Ctx()) and returning
// its error immediately on failure.
func renderDoStatement(w *strings.Builder, e ast.Exec, key string) {
	call := fmt.Sprintf("%s(t.Ctx())", refSelector(e.Func))
	indent := "\t\t"
	if e.Certainty != nil {
		fmt.Fprintf(w, "%sif %s {\n", indent, certaintyExpr(e.Certainty))
		indent += "\t"
	}
	if e.Once {
		fmt.Fprintf(w, "%sif !tangl.OnceDone(t.Ctx(), %s) {\n", indent, strconv.Quote(key))
		fmt.Fprintf(w, "%s\tif err := %s; err != nil {\n", indent, call)
		fmt.Fprintf(w, "%s\t\treturn err\n", indent)
		fmt.Fprintf(w, "%s\t}\n", indent)
		fmt.Fprintf(w, "%s\ttangl.OnceMark(t.Ctx(), %s)\n", indent, strconv.Quote(key))
		fmt.Fprintf(w, "%s}\n", indent)
	} else {
		fmt.Fprintf(w, "%sif err := %s; err != nil {\n", indent, call)
		fmt.Fprintf(w, "%s\treturn err\n", indent)
		fmt.Fprintf(w, "%s}\n", indent)
	}
	if e.Certainty != nil {
		fmt.Fprintln(w, "\t\t}")
	}
}
