//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what writeExtraChecks — writes trailing Evaluate calls concatenating check results
package gen

import (
	"fmt"
	"strings"
)

// writeExtraChecks writes the trailing Evaluate calls for an endpoint
// that checks one or more cases, concatenating every checked case's
// results and returning them alongside a nil error. When returnsResults is
// false (a run-only endpoint with no checks), it simply returns nil.
func writeExtraChecks(w *strings.Builder, checks []string, returnsResults bool) {
	if !returnsResults {
		fmt.Fprintln(w, "\treturn nil")
		return
	}
	fmt.Fprintln(w, "\tvar out []toulmin.EvalResult")
	for i, c := range checks {
		fmt.Fprintf(w, "\tr%d, err := %sGraph.Evaluate(ctx)\n", i, goIdent(c))
		fmt.Fprintln(w, "\tif err != nil {")
		fmt.Fprintln(w, "\t\treturn nil, err")
		fmt.Fprintln(w, "\t}")
		fmt.Fprintf(w, "\tout = append(out, r%d...)\n", i)
	}
	fmt.Fprintln(w, "\treturn out, nil")
}
