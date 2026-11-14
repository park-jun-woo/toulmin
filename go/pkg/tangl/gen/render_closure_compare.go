//ff:func feature=tangl type=codegen control=sequence
//ff:what renderClosureCompare — emits a Compare leaf as a verdict+Fact closure literal
package gen

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureCompare emits a Compare leaf as a func() (bool, any) closure
// literal mirroring pkg/tangl/eval.compileCompare: the left side resolves via
// its operand code (a nil Lhs is the flat Field as a one-step ctx path), the
// right side is the literal Value or a resolved Rhs operand whose error/ok is
// ignored, and the witness Fact carries the Where notation, "<op> <rhs>"
// Expected (the bare op when there is no right side), and the measured left
// value (or a lhs Call's error) as Actual.
func renderClosureCompare(c ast.Compare, env map[string]binderInfo) (string, error) {
	var lc valCode
	whereExpr := strconv.Quote(c.Field)
	if c.Lhs == nil {
		// A flat Field is a one-step path: it resolves to the bound element
		// when it names an in-scope exists binder, else to a ctx key —
		// exactly eval's pathVal on Path{Steps: [Field]}. The Where stays the
		// plain Field notation either way, mirroring eval.compileCompare.
		lc = pathValCode(ast.Path{Steps: []string{c.Field}}, env, "l", true)
	} else {
		var err error
		lc, err = renderOperandVal(c.Lhs, env, "l", true, true)
		if err != nil {
			return "", err
		}
		whereExpr = renderOperandRepr(c.Lhs, env).Expr
	}
	rhsArg := strconv.Quote(c.Value)
	expectedExpr := strconv.Quote(c.Op)
	if c.Value != "" {
		expectedExpr = strconv.Quote(c.Op + " " + c.Value)
	}
	setup := lc.Setup
	if c.Rhs != nil {
		rc, err := renderOperandVal(c.Rhs, env, "r", false, false)
		if err != nil {
			return "", err
		}
		setup += rc.Setup
		rhsArg = rc.V
		rr := renderOperandRepr(c.Rhs, env)
		if rr.Const {
			expectedExpr = strconv.Quote(c.Op + " " + rr.Val)
		} else {
			expectedExpr = strconv.Quote(c.Op+" ") + "+" + rr.Expr
		}
	}
	return "func() (bool, any) {\n" + setup +
		"return tanglCompareVal(" + lc.V + ", " + lc.OK + ", " + strconv.Quote(c.Op) + ", " + rhsArg + "), " +
		"toulmin.Fact{Where: " + whereExpr + ", Expected: " + expectedExpr + ", Actual: " + lc.Actual + "}\n}", nil
}
