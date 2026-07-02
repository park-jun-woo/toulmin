//ff:func feature=tangl type=codegen control=sequence
//ff:what renderInternalOn — writes an unexported handler for an on-event trigger
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderInternalOn writes an unexported handler for an "on <event>"
// trigger: a single Init/Run/Compensate/Commit pass over its Runs cases
// (if any), plus any Checks cases appended as plain Evaluate results.
// Wiring the handler to its actual event source is left as a TODO — tangl
// only declares the event name, not a transport.
func renderInternalOn(w *strings.Builder, in ast.Internal, idx int) {
	fnName := internalFuncName("on", in.Event, idx)
	returnsResults := len(in.Checks) > 0
	sig := "error"
	if returnsResults {
		sig = "([]toulmin.EvalResult, error)"
	}
	fmt.Fprintf(w, "// %s handles the %q event. TODO: wire this to the event source.\n", fnName, in.Event)
	fmt.Fprintf(w, "func %s(ctx toulmin.Context) %s {\n", fnName, sig)
	if len(in.Runs) > 0 {
		fmt.Fprintln(w, "\ttangl.InitCompensation(ctx)")
		writeRunSequence(w, in.Runs, returnsResults)
		fmt.Fprintln(w, "\ttangl.CommitCompensation(ctx)")
	}
	writeExtraChecks(w, in.Checks, returnsResults)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
