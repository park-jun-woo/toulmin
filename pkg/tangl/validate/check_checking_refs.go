//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkCheckingRefs — verifies checking clause targets name an existing case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkCheckingRefs verifies that every node's `checking <case>` clause
// names a case that exists in the document.
func checkCheckingRefs(doc *ast.Document) []error {
	cases := caseNameSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking == "" {
				continue
			}
			if !cases[n.Checking] {
				errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q checking target case %q does not exist", c.Name, n.Name, n.Checking))
			}
		}
	}
	return errs
}
