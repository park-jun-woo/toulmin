//ff:func feature=tangl type=validator control=sequence
//ff:what walkOperand — visit a single condition operand (Path, SpecRef, Call), nil-safe
package validate

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// walkOperand calls fn on operand o when it is non-nil. Operands do not nest,
// so this is a nil-safe single visit that complements walkExpr for the operand
// slots (Compare.Lhs/Rhs, Exists.Source).
func walkOperand(o ast.Operand, fn func(ast.Operand)) {
	if o == nil {
		return
	}
	fn(o)
}
