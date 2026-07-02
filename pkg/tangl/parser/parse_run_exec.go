//ff:func feature=tangl type=parser control=sequence
//ff:what parseRunExec — parse "NAME when NAME" for a run cascade edge
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseRunExec parses run_edge (the "run" keyword already consumed).
func parseRunExec(rest string, line int, path string) (ast.Exec, bool, error) {
	caseName, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected backtick-quoted case name after 'run'")
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected 'when' in run edge")
	}
	node, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected backtick-quoted node name after 'when'")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Exec{}, true, errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return ast.Exec{Kind: ast.RunExec, Case: caseName, Node: node, Line: line}, true, nil
}
