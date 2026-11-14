//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeTickChecks — writes a discarded Evaluate call per ticker check case
package gen

import (
	"fmt"
	"strings"
)

// writeTickChecks writes an Evaluate call per check case in the ticker
// body. Results and errors are both discarded — Internal has no caller to
// report to, and a check case's Evaluate exists for its side-effect-free
// judgment alone.
func writeTickChecks(w *strings.Builder, checks []string) {
	for _, c := range checks {
		fmt.Fprintf(w, "\t\t_, _ = %sGraph.Evaluate(ctx)\n", goIdent(c))
	}
}
