//ff:func feature=tangl type=parser control=sequence
//ff:what parseDoExec — parse "FUNC_REF [once] when NAME [certainty_clause]"
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseDoExec parses do_edge (the "do" keyword already consumed by the caller).
func parseDoExec(rest string, line int, path string) (ast.Exec, bool, error) {
	ref, rest, ok := parseRef(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected function reference after 'do'")
	}
	once := false
	if r2, ok2 := takeKeyword(rest, "once"); ok2 {
		once = true
		rest = r2
	}
	rest, ok = takeKeyword(rest, "when")
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected 'when' in do edge")
	}
	node, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Exec{}, true, errAt(path, line, "expected backtick-quoted node name after 'when'")
	}
	var cert *ast.Certainty
	rest = strings.TrimSpace(rest)
	if rest != "" {
		c, err := parseCertaintyClause(rest, path, line)
		if err != nil {
			return ast.Exec{}, true, err
		}
		cert = c
	}
	return ast.Exec{Kind: ast.DoExec, Func: &ref, Node: node, Once: once, Certainty: cert, Line: line}, true, nil
}
