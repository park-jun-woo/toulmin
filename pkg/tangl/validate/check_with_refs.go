//ff:func feature=tangl type=validator control=iteration dimension=3
//ff:what checkWithRefs — verifies with-clause terms name an existing Definitions entry
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkWithRefs verifies that every `with <term>` argument on a node names a
// term declared in the document's tangl:Definitions section.
func checkWithRefs(doc *ast.Document) []error {
	defs := defNameSet(doc)
	var errs []error
	for _, c := range doc.Cases {
		for _, n := range c.Nodes {
			for _, term := range n.With {
				if !defs[term] {
					errs = append(errs, errAt(doc.Path, n.Line, "case %q: node %q references undefined term %q", c.Name, n.Name, term))
				}
			}
		}
	}
	return errs
}
