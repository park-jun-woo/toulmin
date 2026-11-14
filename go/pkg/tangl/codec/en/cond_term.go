//ff:type feature=tangl type=parser
//ff:what condTerm — one flattened condition-list entry (its and/or prefix + term expression)
package en

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// condTerm is one entry of a flattened condition list: the chaining Prefix
// ("" for the first term, otherwise "and"/"or") and the Term expression to
// render on that list item.
type condTerm struct {
	Prefix string
	Term   ast.Expr
}
