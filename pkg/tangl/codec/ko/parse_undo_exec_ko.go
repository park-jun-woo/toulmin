//ff:func feature=tangl type=parser control=sequence
//ff:what parseUndoExecKo — parse an undo edge body '`X` 면 `F` 로'
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseUndoExecKo parses the body of an undo compensation edge (되돌린다 already
// consumed): the condition node `X` first, then the function ref `F` (the trailing
// 로 is dropped as a particle). An undo edge carries no 한 번만 and no 확신 gate.
func parseUndoExecKo(body string, line int, path string) (ast.Exec, error) {
	node, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected backtick-quoted node name in undo edge")
	}
	rest = skipParticlesKo(rest)
	ref, r2, ok := parseRefKo(rest)
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected function reference in undo edge")
	}
	if strings.TrimSpace(skipParticlesKo(r2)) != "" {
		return ast.Exec{}, errAtKo(path, line, "unexpected trailing text: %q", r2)
	}
	return ast.Exec{Kind: ast.UndoExec, Func: &ref, Node: node, Line: line}, nil
}
