//ff:func feature=tangl type=parser control=sequence
//ff:what parseNodeItem — parse a "`name` is a ... rule [using|checking ...]" node
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseNodeItem parses node_decl: `name` is (a|an) ROLE_EXPR
// ["using" FUNC_REF [with_clause] [qual_clause] | "checking" NAME].
// It returns ok=false (no error) when the item is not a node declaration.
func parseNodeItem(it item, path string) (ast.Node, bool, error) {
	name, rest, ok := takeBacktick(it.Text)
	if !ok {
		return ast.Node{}, false, nil
	}
	rest, ok = takeKeyword(rest, "is")
	if !ok {
		return ast.Node{}, false, nil
	}
	rest, ok = takeArticle(rest)
	if !ok {
		return ast.Node{}, false, nil
	}
	role, rest, ok := parseRoleExpr(rest)
	if !ok {
		return ast.Node{}, false, errAt(path, it.Line, "%q: expected role (general/counter/except rule)", it.Text)
	}
	n := ast.Node{Name: name, Role: role, Line: it.Line}
	rest = strings.TrimSpace(rest)
	if rest == "" {
		return n, true, nil
	}
	if level, ok := takeKeyword(rest, ","); ok {
		lvl, tail, err := parseLevelClause(level, path, it.Line)
		if err != nil {
			return ast.Node{}, true, err
		}
		if strings.TrimSpace(tail) != "" {
			return ast.Node{}, true, errAt(path, it.Line, "unexpected trailing text: %q", tail)
		}
		n.Level = lvl
		return n, true, nil
	}
	if err := applyNodeClause(&n, rest, it.Line, path); err != nil {
		return ast.Node{}, true, err
	}
	return n, true, nil
}
