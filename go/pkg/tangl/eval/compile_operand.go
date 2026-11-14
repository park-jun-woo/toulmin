//ff:func feature=tangl type=engine control=selection
//ff:what compileOperandValue — compiles a value-position Operand into a value closure
package eval

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// compileOperandValue compiles a value-position Operand into a valueFunc.
// Path and SpecRef resolve purely at runtime; a value-position Call is bound to
// res.Value at compile time and an unresolved reference is an error.
func compileOperandValue(o ast.Operand, res Resolver) (valueFunc, error) {
	switch v := o.(type) {
	case ast.Path:
		return func(sc *scope) (any, bool, error) {
			val, ok := pathVal(v, sc)
			return val, ok, nil
		}, nil
	case ast.SpecRef:
		return func(sc *scope) (any, bool, error) {
			val, ok := specVal(v, sc)
			return val, ok, nil
		}, nil
	case ast.Call:
		fn, ok := res.Value(v.Fn.Alias, v.Fn.Name)
		if !ok {
			return nil, fmt.Errorf("eval: unresolved value %s", callRepr(v))
		}
		return func(sc *scope) (any, bool, error) {
			val, err := fn(sc.ctx, sc.specs)
			if err != nil {
				return nil, false, err
			}
			return val, true, nil
		}, nil
	default:
		return nil, fmt.Errorf("eval: unsupported operand node %T", o)
	}
}
