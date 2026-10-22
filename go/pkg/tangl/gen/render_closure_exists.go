//ff:func feature=tangl type=codegen control=sequence pattern=recursive
//ff:what renderClosureExists — emits an Exists quantifier as a tanglAny closure literal
package gen

import (
	"strconv"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// renderClosureExists emits an Exists quantifier as a func() (bool, any)
// closure literal: the source operand resolves to (value, ok) and tanglAny
// iterates it in index order, running the condition closure with the
// depth-indexed element/index variables bound to the binder — matching
// pkg/tangl/eval.compileExists, including the "any of <source> as <binder>
// satisfies" Expected notation and the element count Actual on a false
// verdict.
func renderClosureExists(e ast.Exists, env map[string]binderInfo, depth int) (string, error) {
	srcRepr := renderOperandRepr(e.Source, env)
	sc, err := renderOperandVal(e.Source, env, "s", true, false)
	if err != nil {
		return "", err
	}
	elem := "e" + strconv.Itoa(depth)
	idx := "i" + strconv.Itoa(depth)
	inner := cloneBinderEnv(env, e.Binder, binderInfo{Elem: elem, Idx: idx, Src: srcRepr})
	cond, err := renderClosureExpr(e.Cond, inner, depth+1)
	if err != nil {
		return "", err
	}
	expectedExpr := strconv.Quote("any of ") + "+" + srcRepr.Expr + "+" + strconv.Quote(" as "+e.Binder+" satisfies")
	if srcRepr.Const {
		expectedExpr = strconv.Quote("any of " + srcRepr.Val + " as " + e.Binder + " satisfies")
	}
	return "func() (bool, any) {\n" + sc.Setup +
		"return tanglAny(" + sc.V + ", " + sc.OK + ", " + srcRepr.Expr + ", " + expectedExpr +
		", func(" + elem + " any, " + idx + " int) (bool, any) {\nreturn " + cond + "()\n})\n}", nil
}
