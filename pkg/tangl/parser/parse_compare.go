//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCompare — parse a leaf `field operator value` comparison
package parser

import (
	"strings"

	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// compareOps lists the eight recognized comparison operators.
var compareOps = []string{
	"is not empty", "is empty", "is greater than", "is less than",
	"is at most", "is at least", "is in", "equals",
}

// parseCompare parses "`<field>` <operator> [<value>]".
func parseCompare(text string, line int, path string) (ast.Compare, error) {
	field, rest, ok := takeBacktick(text)
	if !ok {
		return ast.Compare{}, errAt(path, line, "expected backtick-quoted field, got %q", text)
	}
	for _, op := range compareOps {
		if r2, ok2 := takeKeyword(rest, op); ok2 {
			return ast.Compare{Field: field, Op: op, Value: strings.TrimSpace(r2)}, nil
		}
	}
	return ast.Compare{}, errAt(path, line, "unrecognized comparison operator in %q", text)
}
