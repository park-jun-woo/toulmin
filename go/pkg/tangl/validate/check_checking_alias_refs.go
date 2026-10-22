//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkCheckingAliasRefs — verifies a node's checking-clause package alias is declared in See
package validate

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// checkCheckingAliasRefs verifies that every node's `checking <alias>.<case>`
// package alias was declared in the document's tangl:See section.
func checkCheckingAliasRefs(doc *ast.Document) []error {
	aliases := seeAliasSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Checking == nil || n.Checking.Alias == "" {
				continue
			}
			if !aliases[n.Checking.Alias] {
				errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q checks undeclared package alias %q", c.Name, n.Name, n.Checking.Alias))
			}
		}
	}
	return errs
}
