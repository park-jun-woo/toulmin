//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkAttackRefs — verifies don't/do-not target and attacker name registered nodes
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkAttackRefs verifies that every Attack's Target and Attacker name a
// node registered in the same case.
func checkAttackRefs(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		nodes := caseNodeSet(c)
		for _, a := range c.Attacks {
			if !nodes[a.Target] {
				errs = append(errs, errAt(doc.Path, a.Line, "case %q: don't target %q is not a registered node", c.Name, a.Target))
			}
			if !nodes[a.Attacker] {
				errs = append(errs, errAt(doc.Path, a.Line, "case %q: don't attacker %q is not a registered node", c.Name, a.Attacker))
			}
		}
	}
	return errs
}
