//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what appendChildren — flatten child expressions in order onto a record buffer, summing node counts
package packet

import (
	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// appendChildren flattens each child of terms in order via emitExpr, appending
// their records to out and adding their node counts to count. It is the shared
// recursion step for Logic/Not/Either records (§design 4).
func appendChildren(out []word.Word, count int, terms []ast.Expr, labels *[]string) ([]word.Word, int, error) {
	for _, t := range terms {
		cw, c, err := emitExpr(t, labels)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, cw...)
		count += c
	}
	return out, count, nil
}
