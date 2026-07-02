//ff:func feature=tangl type=parser control=iteration dimension=1 pattern=recursive
//ff:what parseCondList — parse a sibling condition list, chained by and/or
package parser

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// parseCondList parses a sibling list of condition items: the first item is
// the seed expression, and each following item must be prefixed with "and"
// or "or", chaining onto the running result.
func parseCondList(items []item, path string) (ast.Expr, error) {
	if len(items) == 0 {
		return nil, errAt(path, 0, "expected at least one condition")
	}
	result, err := parseCondTerm(items[0], path)
	if err != nil {
		return nil, err
	}
	for _, it := range items[1:] {
		op, text, ok := takeCondPrefix(it.Text)
		if !ok {
			return nil, errAt(path, it.Line, "expected 'and'/'or' prefix, got %q", it.Text)
		}
		term, err := parseCondTerm(item{Text: text, Line: it.Line, Children: it.Children}, path)
		if err != nil {
			return nil, err
		}
		result = mergeLogic(result, op, term)
	}
	return result, nil
}
