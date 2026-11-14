//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what registerCaseSupersedes — writes one Supersedes statement per supersede edge
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// registerCaseSupersedes writes one superseder.Supersedes(target) statement
// per "`x` supersedes `y`" edge, in document order.
func registerCaseSupersedes(w *strings.Builder, c ast.Case, nodes map[string]nodeInfo) {
	seen := map[[2]string]bool{}
	for _, s := range c.Supersedes {
		key := [2]string{s.Superseder, s.Superseded}
		if seen[key] {
			continue
		}
		seen[key] = true
		superseder := nodes[s.Superseder].Var
		target := nodes[s.Superseded].Var
		fmt.Fprintf(w, "\t%s.Supersedes(%s)\n", superseder, target)
	}
	fmt.Fprintln(w)
}
