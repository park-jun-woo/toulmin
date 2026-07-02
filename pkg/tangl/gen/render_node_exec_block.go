//ff:func feature=tangl type=codegen control=sequence
//ff:what renderNodeExecBlock — writes one node's do/undo/run execution attachment
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderNodeExecBlock writes one node's execution attachment: a RunOn
// handler composing its do/undo edges (if any), and/or a .Run(subgraph)
// attachment (if any) — chained onto the same statement when both are
// present.
func renderNodeExecBlock(w *strings.Builder, subject, caseName string, ni nodeInfo, execs []ast.Exec) error {
	runCase, err := runTarget(caseName, ni.Node.Name, execs)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "\t%s", ni.Var)
	if hasDoOrUndo(execs) {
		fmt.Fprintln(w, ".RunOn(func(self toulmin.TraceEntry, t toulmin.Trace) error {")
		writeNodeExecs(w, subject, caseName, ni.Node.Name, execs)
		fmt.Fprintln(w, "\t\treturn nil")
		fmt.Fprint(w, "\t})")
	}
	if runCase != "" {
		fmt.Fprintf(w, ".Run(%sGraph)", goIdent(runCase))
	}
	fmt.Fprintln(w)
	return nil
}
