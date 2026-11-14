//ff:func feature=tangl type=parser control=sequence
//ff:what parseDoExecKo — parse a do edge body '`X` [확신 게이트] 면 `F` [한 번만]'
package ko

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseDoExecKo parses the body of a do edge (실행한다 already consumed): the
// condition node `X` comes first, then an optional 확신 gate, then the function ref
// `F`, then an optional 한 번만 (once). The trailing 한 번만 is stripped first, the
// node backtick is read from the front, and the text between the node and the
// function backtick is the condition — a 확신 gate (parseCertaintyKo) or a bare 면.
func parseDoExecKo(body string, line int, path string) (ast.Exec, error) {
	once := false
	if b, ok := takeSuffixKo(body, "한 번만"); ok {
		once = true
		body = b
	}
	node, rest, ok := takeBacktick(body)
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected backtick-quoted node name in do edge")
	}
	fi := strings.Index(rest, "`")
	if fi < 0 {
		return ast.Exec{}, errAtKo(path, line, "expected function reference in do edge")
	}
	condText := rest[:fi]
	ref, r2, ok := parseRefKo(rest[fi:])
	if !ok {
		return ast.Exec{}, errAtKo(path, line, "expected function reference after node in do edge")
	}
	if strings.TrimSpace(skipParticlesKo(r2)) != "" {
		return ast.Exec{}, errAtKo(path, line, "unexpected trailing text: %q", r2)
	}
	var cert *ast.Certainty
	if strings.Contains(condText, "확신") {
		c, err := parseCertaintyKo(condText, line, path)
		if err != nil {
			return ast.Exec{}, err
		}
		cert = c
	} else if strings.TrimSpace(skipParticlesKo(condText)) != "" {
		return ast.Exec{}, errAtKo(path, line, "unexpected condition text in do edge: %q", condText)
	}
	return ast.Exec{Kind: ast.DoExec, Func: &ref, Node: node, Once: once, Certainty: cert, Line: line}, nil
}
