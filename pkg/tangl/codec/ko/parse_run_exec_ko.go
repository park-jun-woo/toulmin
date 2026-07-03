//ff:func feature=tangl type=parser control=sequence
//ff:what parseRunExecKo — parse a run cascade edge body '`X` 면 `C`'
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseRunExecKo parses the body of a run cascade edge (진행한다 already consumed):
// the condition node `X` first, then the target case name `C` (a bare string, not a
// function ref). It is the case-scoped run (with a 면 condition); the endpoint/
// internal body run (no condition) is parsed by parseRunCheckKo instead.
func parseRunExecKo(body string, line int, path string) (ast.Exec, error) {
	node, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected backtick-quoted node name in run edge")
	}
	rest = skipParticlesKo(rest)
	caseName, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected backtick-quoted case name in run edge")
	}
	if strings.TrimSpace(skipParticlesKo(rest)) != "" {
		return ast.Exec{}, errAtKo(path, line, "unexpected trailing text: %q", rest)
	}
	return ast.Exec{Kind: ast.RunExec, Case: caseName, Node: node, Line: line}, nil
}
