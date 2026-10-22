//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkSupersedeRefs — verifies supersede superseder and superseded name registered nodes
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkSupersedeRefs verifies that every Supersede's Superseder and
// Superseded name a node registered in the same case.
func checkSupersedeRefs(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		nodes := caseNodeSet(c)
		for _, s := range c.Supersedes {
			if !nodes[s.Superseder] {
				errs = append(errs, errAt(doc.Path, s.Line, "case %q: supersedes superseder %q is not a registered node", c.Name, s.Superseder))
			}
			if !nodes[s.Superseded] {
				errs = append(errs, errAt(doc.Path, s.Line, "case %q: supersedes target %q is not a registered node", c.Name, s.Superseded))
			}
		}
	}
	return errs
}
