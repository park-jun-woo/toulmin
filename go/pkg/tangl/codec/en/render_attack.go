//ff:func feature=tangl type=parser control=sequence
//ff:what renderAttack — render a defeat edge (`don't `target` when `attacker``)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderAttack renders a defeat edge back to "don't `<target>` when `<attacker>`",
// the canonical form of the parser's don't / do not edge.
func renderAttack(a ast.Attack) string {
	return "don't `" + a.Target + "` when `" + a.Attacker + "`"
}
