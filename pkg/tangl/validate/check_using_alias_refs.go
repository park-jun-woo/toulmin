//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkUsingAliasRefs — verifies a node's using-clause package alias is declared in See
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkUsingAliasRefs verifies that every node's `using <alias>.<name>`
// package alias was declared in the document's tangl:See section.
func checkUsingAliasRefs(doc *ast.Document) []error {
	aliases := seeAliasSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			if n.Using == nil || n.Using.Alias == "" {
				continue
			}
			if !aliases[n.Using.Alias] {
				errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q uses undeclared package alias %q", c.Name, n.Name, n.Using.Alias))
			}
		}
	}
	return errs
}
