//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkExecNodeRefs — verifies do/undo/run "when" nodes are registered in the same case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkExecNodeRefs verifies that every Exec's trigger Node names a node
// registered in the same case, for all of do/undo/run.
func checkExecNodeRefs(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		nodes := caseNodeSet(c)
		for _, e := range c.Execs {
			if !nodes[e.Node] {
				errs = append(errs, errAt(doc.Path, e.Line, "case %q: exec 'when' node %q is not a registered node", c.Name, e.Node))
			}
		}
	}
	return errs
}
