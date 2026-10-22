//ff:func feature=tangl type=codegen control=sequence
//ff:what renderEndpointRun — writes a Provides run endpoint's compensation-wrapped function
package gen

import (
	"fmt"
	"strings"
)

// renderEndpointRun writes a Provides "run" endpoint: Required() guards
// its fields, then the spec's compensation wrapper (InitCompensation,
// each case's Run in document order, Compensate/Review on error,
// CommitCompensation on success). If extraChecks is non-empty the
// function also Evaluates those cases afterward and returns their
// combined results alongside the error.
func renderEndpointRun(w *strings.Builder, fnName string, fields, runs, extraChecks []string) {
	returnsResults := len(extraChecks) > 0
	sig := "error"
	if returnsResults {
		sig = "([]toulmin.EvalResult, error)"
	}
	fmt.Fprintf(w, "func %s(ctx toulmin.Context) %s {\n", fnName, sig)
	writeRequiredGuard(w, fields, returnsResults)
	fmt.Fprintln(w, "\ttangl.InitCompensation(ctx)")
	writeRunSequence(w, runs, returnsResults)
	fmt.Fprintln(w, "\ttangl.CommitCompensation(ctx)")
	writeExtraChecks(w, extraChecks, returnsResults)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
