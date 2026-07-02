//ff:func feature=tangl type=validator control=selection
//ff:what armUndoExec — arms a node on do, or reports missing-do on undo
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// armUndoExec updates armed for a do exec, or (for an undo exec) returns an
// error when the node was never armed by a preceding do (spec §undo).
func armUndoExec(path string, caseName string, armed map[string]bool, e ast.Exec) error {
	switch e.Kind {
	case ast.DoExec:
		armed[e.Node] = true
	case ast.UndoExec:
		if !armed[e.Node] {
			return errAt(path, e.Line, "case %q: undo when %q has no preceding do on that node", caseName, e.Node)
		}
	}
	return nil
}
