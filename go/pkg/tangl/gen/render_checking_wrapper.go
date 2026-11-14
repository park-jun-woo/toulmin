//ff:func feature=tangl type=codegen control=sequence
//ff:what renderCheckingWrapper — writes a pure wrapper that Evaluates the target case (errors propagate, never a silent -1)
package gen

import (
	"fmt"
	"strings"
)

// renderCheckingWrapper writes a pure
// func(toulmin.Context, toulmin.Specs) (bool, any) that Evaluates the
// target case's graph and composes its warrant results into a single
// (active, results) verdict via tanglCaseActive — "checking" never fires
// do/run since Evaluate has no side effects. A failed Evaluate is not folded
// into false (−1, "does not apply"): the rule signature has no error slot, so
// the wrapper panics with the error and the engine's safeCall turns it into
// the calling graph's Evaluate/Run error (failure is never silent).
func renderCheckingWrapper(w *strings.Builder, fnName, targetCase string) {
	fmt.Fprintf(w, "func %s(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {\n", fnName)
	fmt.Fprintf(w, "\tresults, err := %sGraph.Evaluate(ctx)\n", goIdent(targetCase))
	fmt.Fprintln(w, "\tif err != nil {")
	fmt.Fprintln(w, "\t\tpanic(err) // surfaced by the engine as this rule's Evaluate/Run error")
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "\treturn tanglCaseActive(results), results")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
