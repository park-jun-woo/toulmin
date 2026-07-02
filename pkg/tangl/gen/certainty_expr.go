//ff:func feature=tangl type=codegen control=selection
//ff:what certaintyExpr — converts a certainty clause into a self.Verdict comparison
package gen

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// certaintyExpr converts an "if <op> N% certain" clause into a
// self.Verdict comparison, using the spec's verdict/percent inverse:
// verdict = N/50 - 1.
func certaintyExpr(c *ast.Certainty) string {
	threshold := formatFloat(float64(c.Percent)/50.0 - 1.0)
	switch c.Op {
	case "above":
		return fmt.Sprintf("self.Verdict > %s", threshold)
	case "less than":
		return fmt.Sprintf("self.Verdict < %s", threshold)
	case "at most":
		return fmt.Sprintf("self.Verdict <= %s", threshold)
	default: // "at least"
		return fmt.Sprintf("self.Verdict >= %s", threshold)
	}
}
