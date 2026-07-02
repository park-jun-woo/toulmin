//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what registerCaseAttacks — writes one Attacks statement per defeat edge
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// registerCaseAttacks writes one attacker.Attacks(target) statement per
// "don't <target> when <attacker>" defeat edge, in document order.
func registerCaseAttacks(w *strings.Builder, c ast.Case, nodes map[string]nodeInfo) {
	for _, a := range c.Attacks {
		attacker := nodes[a.Attacker].Var
		target := nodes[a.Target].Var
		fmt.Fprintf(w, "\t%s.Attacks(%s)\n", attacker, target)
	}
	fmt.Fprintln(w)
}
