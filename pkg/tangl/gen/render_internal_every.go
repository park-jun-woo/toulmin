//ff:func feature=tangl type=codegen control=sequence
//ff:what renderInternalEvery — writes an unexported ticker-driven runner
package gen

import (
	"fmt"
	"strings"
	"time"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderInternalEvery writes an unexported ticker-driven runner for an
// "every <interval> [until <case>]" trigger: one Init/Run/Compensate/
// Commit compensation cycle per tick, sharing ctx across the whole loop
// (the "once" guard's lifetime), first checking the "until" case's
// Evaluate result (if any) to end the loop.
func renderInternalEvery(w *strings.Builder, in ast.Internal, idx int) {
	fnName := tickFuncName(in, idx)
	d, ok := normalizeInterval(in.Interval)
	comment := fmt.Sprintf("interval: %q", in.Interval)
	if !ok {
		d = 24 * time.Hour
		comment = fmt.Sprintf("interval %q is not a supported duration/clock schedule, defaulting to 24h", in.Interval)
	}
	fmt.Fprintf(w, "func %s(ctx toulmin.Context) {\n", fnName)
	fmt.Fprintf(w, "\tticker := time.NewTicker(time.Duration(%d)) // %s\n", int64(d), comment)
	fmt.Fprintln(w, "\tdefer ticker.Stop()")
	fmt.Fprintln(w, "\tfor range ticker.C {")
	if in.Until != "" {
		fmt.Fprintf(w, "\t\tresults, err := %sGraph.Evaluate(ctx)\n", goIdent(in.Until))
		fmt.Fprintln(w, "\t\tif err == nil && tanglCaseActive(results) {")
		fmt.Fprintln(w, "\t\t\treturn")
		fmt.Fprintln(w, "\t\t}")
	}
	fmt.Fprintln(w, "\t\ttangl.InitCompensation(ctx)")
	writeTickRunSequence(w, in.Runs)
	writeTickChecks(w, in.Checks)
	fmt.Fprintln(w, "\t\ttangl.CommitCompensation(ctx)")
	fmt.Fprintln(w, "\t}")
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
