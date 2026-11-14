//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what parseCondListKo — parse a sibling condition list, chained by 그리고/또는
package ko

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/codec/scan"
)

// parseCondListKo parses a sibling list of condition items: the first item is the
// seed expression, and each following item must be prefixed with 그리고 or 또는,
// chaining onto the running result as an ast.Logic ("and"/"or"). It mirrors
// pkg/tangl/parser.parseCondList.
func parseCondListKo(items []scan.Item, path string) (ast.Expr, error) {
	if len(items) == 0 {
		return nil, errAtKo(path, 0, "expected at least one condition")
	}
	result, err := parseCondTermKo(items[0].Text, items[0].Children, items[0].No, path)
	if err != nil {
		return nil, err
	}
	for _, it := range items[1:] {
		op, text, ok := takeCondPrefixKo(it.Text)
		if !ok {
			return nil, errAtKo(path, it.No, "expected '그리고'/'또는' prefix, got %q", it.Text)
		}
		term, err := parseCondTermKo(text, it.Children, it.No, path)
		if err != nil {
			return nil, err
		}
		result = ast.Logic{Op: op, Terms: []ast.Expr{result, term}}
	}
	return result, nil
}
