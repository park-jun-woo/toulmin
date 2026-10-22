//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what parseChildren — rebuild n child subtrees in prefix order, summing their node counts
package packet

import (
	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// parseChildren rebuilds n child subtrees in order via parseExpr, advancing *pos
// and returning the children with the total node count they consumed. It is the
// shared recursion step for Logic/Not/Either records.
func parseChildren(recs []core.Word, pos *int, labels []string, n int) ([]ast.Expr, int, error) {
	terms := make([]ast.Expr, 0, n)
	count := 0
	for i := 0; i < n; i++ {
		child, c, err := parseExpr(recs, pos, labels)
		if err != nil {
			return nil, 0, err
		}
		terms = append(terms, child)
		count += c
	}
	return terms, count, nil
}
