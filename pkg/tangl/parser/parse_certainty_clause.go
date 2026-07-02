//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCertaintyClause — parse "if THRESHOLD_OP INT % certain"
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// certaintyOps lists the four recognized threshold operators.
var certaintyOps = []string{"at least", "above", "less than", "at most"}

// parseCertaintyClause parses "if <op> <N>% certain".
func parseCertaintyClause(s string, path string, line int) (*ast.Certainty, error) {
	rest, ok := takeKeyword(s, "if")
	if !ok {
		return nil, errAt(path, line, "expected 'if' certainty clause, got %q", s)
	}
	var op string
	for _, o := range certaintyOps {
		if r2, ok2 := takeKeyword(rest, o); ok2 {
			op, rest = o, r2
			break
		}
	}
	if op == "" {
		return nil, errAt(path, line, "expected certainty operator, got %q", rest)
	}
	rest = trimStart(rest)
	i := 0
	for i < len(rest) && rest[i] >= '0' && rest[i] <= '9' {
		i++
	}
	if i == 0 {
		return nil, errAt(path, line, "expected integer percent, got %q", rest)
	}
	pct, rest, err := atoiRest(rest, i, path, line)
	if err != nil {
		return nil, err
	}
	rest, ok = takeKeyword(rest, "%")
	if !ok {
		return nil, errAt(path, line, "expected '%%' after percent, got %q", rest)
	}
	rest, ok = takeKeyword(rest, "certain")
	if !ok {
		return nil, errAt(path, line, "expected 'certain', got %q", rest)
	}
	if strings.TrimSpace(rest) != "" {
		return nil, errAt(path, line, "unexpected trailing text: %q", rest)
	}
	return &ast.Certainty{Op: op, Percent: pct}, nil
}
