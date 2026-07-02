//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what unguardedDoWarnings — warns about reachable do edges with no once guard
package validate

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// unguardedDoWarnings warns about every `do` edge, in a case reached from
// doc's tangl:Internal, that has no `once` guard (spec §틱 멱등성 요구).
func unguardedDoWarnings(doc *ast.Document, reached map[string]bool) []string {
	var warnings []string
	for _, c := range doc.Cases {
		if !reached[c.Name] {
			continue
		}
		for _, e := range c.Execs {
			if e.Kind != ast.DoExec || e.Once {
				continue
			}
			warnings = append(warnings, fmt.Sprintf(
				"%s:%d: warning: case %q do %s has no once guard and is reachable from tangl:Internal",
				doc.Path, e.Line, c.Name, refString(e.Func)))
		}
	}
	return warnings
}
