//ff:func feature=tangl type=codegen control=sequence
//ff:what renderCheckingWrapper — writes a pure wrapper that Evaluates the target case
package gen

import (
	"fmt"
	"strings"
)

// renderCheckingWrapper writes a pure
// func(toulmin.Context, toulmin.Specs) (bool, any) that Evaluates the
// target case's graph and composes its warrant results into a single
// (active, results) verdict via tanglCaseActive — "checking" never fires
// do/run since Evaluate has no side effects.
func renderCheckingWrapper(w *strings.Builder, fnName, targetCase string) {
	fmt.Fprintf(w, "func %s(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {\n", fnName)
	fmt.Fprintf(w, "\tresults, err := %sGraph.Evaluate(ctx)\n", goIdent(targetCase))
	fmt.Fprintln(w, "\tif err != nil {")
	fmt.Fprintln(w, "\t\treturn false, err")
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "\treturn tanglCaseActive(results), results")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
