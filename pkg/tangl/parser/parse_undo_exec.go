//ff:func feature=tangl type=parser control=sequence
//ff:what parseUndoExec — parse "FUNC_REF when NAME" for an undo compensation edge
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseUndoExec parses undo_edge (the "undo" keyword already consumed).
func parseUndoExec(rest string, line int, path string) (ast.Exec, bool, error) {
	ref, rest, ok := parseRef(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected function reference after 'undo'")
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected 'when' in undo edge")
	}
	node, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected backtick-quoted node name after 'when'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Exec{}, true, errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return ast.Exec{Kind: ast.UndoExec, Func: &ref, Node: node, Line: line}, true, nil
}
