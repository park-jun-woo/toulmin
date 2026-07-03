//ff:func feature=tangl type=parser control=sequence pattern=recursive
//ff:what flattenLogic — flatten a left-associative and/or Logic chain into a prefixed term list
package en

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// flattenLogic reverses the parser's left-associative condition-list fold: a
// Logic{op,[L,R]} node contributes R as a term prefixed by op, chained onto the
// flattening of L (which may itself be a Logic). A non-Logic expression is a
// single unprefixed term. The result re-parses, via parseCondList's fold, to the
// original tree.
func flattenLogic(e ast.Expr) []condTerm {
	if lg, ok := e.(ast.Logic); ok && len(lg.Terms) == 2 {
		return append(flattenLogic(lg.Terms[0]), condTerm{Prefix: lg.Op, Term: lg.Terms[1]})
	}
	return []condTerm{{Prefix: "", Term: e}}
}
