//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what flattenLogicKo — flatten a left-associative and/or Logic chain into 그리고/또는-prefixed terms
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// flattenLogicKo reverses parseCondListKo's left-associative fold: a Logic{op,[L,R]}
// node contributes R as a term prefixed by the KO chaining word (그리고 for "and",
// 또는 for "or"), chained onto the flattening of L (which may itself be a Logic). A
// non-Logic expression is a single unprefixed term. The result re-parses, via the
// sibling-list fold, to the original tree.
func flattenLogicKo(e ast.Expr) []condTermKo {
	if lg, ok := e.(ast.Logic); ok && len(lg.Terms) == 2 {
		word := "그리고"
		if lg.Op == "or" {
			word = "또는"
		}
		return append(flattenLogicKo(lg.Terms[0]), condTermKo{Prefix: word, Term: lg.Terms[1]})
	}
	return []condTermKo{{Prefix: "", Term: e}}
}
