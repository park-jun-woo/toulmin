//ff:func feature=tangl type=codegen control=iteration dimension=2
//ff:what collectCheckingTargets — returns the set of case names referenced by any node's checking clause
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// collectCheckingTargets returns the set of case names referenced by any
// node's "checking" clause across the whole document.
func collectCheckingTargets(doc *ast.Document) map[string]bool {
	targets := make(map[string]bool)
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking != "" {
				targets[n.Checking] = true
			}
		}
	}
	return targets
}
