//ff:func feature=tangl type=parser control=iteration dimension=1
//ff:what parseCompare — parse a leaf `field operator value` comparison
package parser

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// compareOps lists the eight recognized comparison operators.
var compareOps = []string{
	"is not empty", "is empty", "is greater than", "is less than",
	"is at most", "is at least", "is in", "equals",
}

// parseCompare parses "<operand> <operator> [<value> | <operand>]". A
// single-step backtick left side is demoted to Field (Lhs nil) so pre-existing
// documents keep a byte-identical AST; multi-step paths, spec references, and
// host calls ride in Lhs. A right side that begins with an operand introducer
// (backtick, "spec", or "result") is parsed into Rhs; otherwise it stays a
// literal Value string.
func parseCompare(text string, line int, path string) (ast.Compare, error) {
	lhs, rest, ok, err := parseOperand(text, line, path)
	if err != nil {
		return ast.Compare{}, err
	}
	if !ok {
		return ast.Compare{}, errAt(path, line, "expected backtick-quoted field, got %q", text)
	}
	cmp := ast.Compare{}
	if p, isPath := lhs.(ast.Path); isPath && len(p.Steps) == 1 {
		cmp.Field = p.Steps[0]
	} else {
		cmp.Lhs = lhs
	}
	for _, op := range compareOps {
		r2, matched := takeKeyword(rest, op)
		if !matched {
			continue
		}
		cmp.Op = op
		return withCompareRhs(cmp, r2, line, path)
	}
	return ast.Compare{}, errAt(path, line, "unrecognized comparison operator in %q", text)
}
