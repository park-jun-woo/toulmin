//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what renderTermKo — render one condition-list term (compare / negation / 다음 중 either) as list-marker lines
package ko

import (
	"fmt"
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// renderTermKo renders a single condition-list term at depth, prefixed by its
// chaining word (그리고/또는, or none). A leaf Compare or negated Compare is one line;
// an Either group emits a "다음 중" head plus its inner condition list at depth+1. The
// 2-space indentation and `- ` marker are relative; the core renderer shifts the
// whole subtree to the section base.
func renderTermKo(prefix string, term ast.Expr, depth int) ([]string, error) {
	head := strings.Repeat("  ", depth) + "- "
	if prefix != "" {
		head += prefix + " "
	}
	if pred, ok, err := renderCondPredKo(term); err != nil {
		return nil, err
	} else if ok {
		return []string{head + pred}, nil
	}
	if either, ok := term.(ast.Either); ok {
		inner, err := renderCondListKo(either.Terms[0], depth+1)
		if err != nil {
			return nil, err
		}
		return append([]string{head + "다음 중"}, inner...), nil
	}
	return nil, fmt.Errorf("ko: unsupported condition term %T", term)
}
