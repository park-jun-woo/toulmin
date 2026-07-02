//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCaseItem — parse one "in case of `name`" case block
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// parseCaseItem parses "in case of `<name>`" and dispatches its children to
// requires, node registrations, attack edges, or exec edges, preserving the
// document declaration order of exec edges (§실행 순서 결정론).
func parseCaseItem(it item, path string) (ast.Case, error) {
	rest, ok := takeKeyword(it.Text, "in case of")
	if !ok {
		return ast.Case{}, errAt(path, it.Line, "expected 'in case of `name`', got %q", it.Text)
	}
	name, rest, ok := takeBacktick(rest)
	if !ok {
		return ast.Case{}, errAt(path, it.Line, "expected backtick-quoted case name")
	}
	if strings.TrimSpace(rest) != "" {
		return ast.Case{}, errAt(path, it.Line, "unexpected trailing text: %q", rest)
	}
	c := ast.Case{Name: name, Line: it.Line}
	for _, child := range it.Children {
		if err := applyCaseChild(&c, child, path); err != nil {
			return ast.Case{}, err
		}
	}
	return c, nil
}
