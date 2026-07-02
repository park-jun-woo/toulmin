//ff:func feature=tangl type=codegen control=iteration dimension=1
//ff:what hasDoOrUndo — reports whether execs contains any do or undo edge
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// hasDoOrUndo reports whether execs contains any do or undo edge (as
// opposed to only a run edge), which decides whether a node needs a
// RunOn handler at all.
func hasDoOrUndo(execs []ast.Exec) bool {
	for _, e := range execs {
		if e.Kind == ast.DoExec || e.Kind == ast.UndoExec {
			return true
		}
	}
	return false
}
