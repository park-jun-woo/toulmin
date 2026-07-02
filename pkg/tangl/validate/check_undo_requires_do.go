//ff:func feature=tangl type=validator control=iteration dimension=2
//ff:what checkUndoRequiresDo — verifies an undo edge's node has a preceding do in the same case
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// checkUndoRequiresDo verifies that every `undo <func> when <node>` edge is
// preceded, in document order within the same case, by at least one
// `do <func> when <node>` edge on that same node (spec §undo).
func checkUndoRequiresDo(doc *ast.Document) []error {
	var errs []error
	for _, c := range doc.Cases {
		armed := make(map[string]bool)
		for _, e := range c.Execs {
			if err := armUndoExec(doc.Path, c.Name, armed, e); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
