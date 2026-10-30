//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what registerCaseAttacks — writes one Attacks statement per defeat edge (with the explicit force when set)
package gen

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// registerCaseAttacks writes one attacker.Attacks(target) statement per defeat
// edge, in document order. An explicit force (`attacks ... with force F` /
// `F의 세기로 공격한다`) becomes attacker.Attacks(target, F), even for F = 1; the
// omitted form keeps the engine default (force 1).
func registerCaseAttacks(w *strings.Builder, c ast.Case, nodes map[string]nodeInfo) {
	for _, a := range c.Attacks {
		attacker := nodes[a.Attacker].Var
		target := nodes[a.Target].Var
		if a.Force != nil {
			fmt.Fprintf(w, "\t%s.Attacks(%s, %s)\n", attacker, target, formatFloat(*a.Force))
			continue
		}
		fmt.Fprintf(w, "\t%s.Attacks(%s)\n", attacker, target)
	}
	fmt.Fprintln(w)
}
