//ff:type feature=tangl type=parser
//ff:what condTermKo — one flattened condition-list entry (its 그리고/또는 prefix + term expression)
package ko

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// condTermKo is one entry of a flattened condition list: the chaining Prefix
// ("" for the first term, otherwise 그리고/또는) and the Term expression to render
// on that list item.
type condTermKo struct {
	Prefix string
	Term   ast.Expr
}
