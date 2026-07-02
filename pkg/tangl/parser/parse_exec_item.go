//ff:func feature=tangl type=parser control=sequence
//ff:what parseExecItem — dispatch a case-body item to do/undo/run edge parsing
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseExecItem dispatches exec_decl to do_edge, undo_edge, or run_edge. A
// "do" immediately followed by "not" is a don't-edge, not a do-edge, and is
// left for parseAttackItem to handle (ok=false here).
func parseExecItem(it item, path string) (ast.Exec, bool, error) {
	if rest, ok := takeKeyword(it.Text, "do"); ok {
		if _, isNot := takeKeyword(rest, "not"); isNot {
			return ast.Exec{}, false, nil
		}
		return parseDoExec(rest, it.Line, path)
	}
	if rest, ok := takeKeyword(it.Text, "undo"); ok {
		return parseUndoExec(rest, it.Line, path)
	}
	if rest, ok := takeKeyword(it.Text, "run"); ok {
		return parseRunExec(rest, it.Line, path)
	}
	return ast.Exec{}, false, nil
}
