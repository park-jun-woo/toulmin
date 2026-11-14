//ff:func feature=tangl type=codegen control=sequence
//ff:what renderEndpointCheck — writes a Provides check-only endpoint's pure Evaluate function
package gen

import (
	"fmt"
	"strings"
)

// renderEndpointCheck writes a Provides "check"-only endpoint: Required()
// guards its fields, then each checked case's Evaluate results are
// concatenated and returned — no side effects, no compensation.
func renderEndpointCheck(w *strings.Builder, fnName string, fields, checks []string) {
	fmt.Fprintf(w, "func %s(ctx toulmin.Context) ([]toulmin.EvalResult, error) {\n", fnName)
	writeRequiredGuard(w, fields, true)
	writeExtraChecks(w, checks, true)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
