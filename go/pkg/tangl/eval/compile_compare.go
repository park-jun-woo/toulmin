//ff:func feature=tangl type=engine control=sequence
//ff:what compileCompare — compiles a Compare leaf into a verdict+Fact closure
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// compileCompare compiles a Compare leaf. The left side is Compare.Field
// (single-step ctx path) when Lhs is nil, else the Lhs operand; the right side
// is the literal Compare.Value when Rhs is nil, else the Rhs operand. The
// returned closure yields the verdict and a toulmin.Fact witness whose Actual
// carries the resolved lhs value (nil for no-value, or the error of a failed
// value Call).
func compileCompare(c ast.Compare, res Resolver) (evalFunc, error) {
	var lhs valueFunc
	var err error
	if c.Lhs != nil {
		lhs, err = compileOperandValue(c.Lhs, res)
		if err != nil {
			return nil, err
		}
	} else {
		field := c.Field
		lhs = func(sc *scope) (any, bool, error) {
			v, ok := pathVal(ast.Path{Steps: []string{field}}, sc)
			return v, ok, nil
		}
	}
	var rhs valueFunc
	if c.Rhs != nil {
		rhs, err = compileOperandValue(c.Rhs, res)
		if err != nil {
			return nil, err
		}
	}
	op := c.Op
	return func(sc *scope) (bool, any) {
		lv, lok, lerr := lhs(sc)
		var rv any
		var rrepr string
		if rhs == nil {
			rv = c.Value
			rrepr = c.Value
		} else {
			rvv, _, _ := rhs(sc)
			rv = rvv
			rrepr = reprOperand(c.Rhs, sc)
		}
		b := compareVals(lv, lok, op, rv)
		where := c.Field
		if c.Lhs != nil {
			where = reprOperand(c.Lhs, sc)
		}
		var actual any
		if lerr != nil {
			actual = lerr
		} else if lok {
			actual = lv
		}
		expected := op
		if rrepr != "" {
			expected = op + " " + rrepr
		}
		return b, toulmin.Fact{Where: where, Expected: expected, Actual: actual}
	}, nil
}
