//ff:type feature=tangl type=model
//ff:what binding — one in-scope exists binder: its source operand, current index, bound value
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// binding is the state of one in-scope exists binder during evaluation: the
// Source operand it iterates, the current element index, and the bound value.
type binding struct {
	source ast.Operand
	index  int
	value  any
}
