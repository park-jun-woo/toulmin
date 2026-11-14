//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what certaintyGateWarnings — warns about do edges whose certainty gate never passes or always passes
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// certaintyGateWarnings warns about every do edge whose `if <op> N% certain`
// gate is out of the reachable range (50%, 100%] of a firing handler
// (certaintyGateReach): a gate that can never pass makes the edge dead code,
// and one that always passes is a no-op the author probably did not mean.
func certaintyGateWarnings(doc *ast.Document) []string {
	var warnings []string
	for _, c := range doc.Cases {
		for _, e := range c.Execs {
			if e.Certainty == nil {
				continue
			}
			never, always := certaintyGateReach(*e.Certainty)
			if never {
				warnings = append(warnings, fmt.Sprintf(
					"%s:%d: warning: case %q gate 'if %s %d%% certain' never passes: a handler fires only above 50%% certain",
					doc.Path, e.Line, c.Name, e.Certainty.Op, e.Certainty.Percent))
			} else if always {
				warnings = append(warnings, fmt.Sprintf(
					"%s:%d: warning: case %q gate 'if %s %d%% certain' always passes: a handler fires only above 50%% certain",
					doc.Path, e.Line, c.Name, e.Certainty.Op, e.Certainty.Percent))
			}
		}
	}
	return warnings
}
