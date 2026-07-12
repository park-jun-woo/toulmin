//ff:func feature=tangl type=codegen control=selection
//ff:what callValCode — emits the code resolving a value-position host Call
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// callValCode emits the code resolving a value-position host Call
// (`result of f`, Go signature func(toulmin.Context, toulmin.Specs)
// (any, error)), mirroring pkg/tangl/eval.compileOperandValue: an error makes
// the value a no-value (nil, ok=false). In lhs position (needActual) the
// witness Fact's Actual carries the error value itself; in rhs position the
// error is swallowed and the value is nil; in exists-source position only the
// ok flag is added.
func callValCode(c ast.Call, prefix string, needOK, needActual bool) valCode {
	v := prefix + "v"
	errV := prefix + "err"
	ok := prefix + "ok"
	actual := prefix + "actual"
	fn := callFnSelector(c)
	head := v + ", " + errV + " := " + fn + "(ctx, specs)\n"
	switch {
	case needActual:
		setup := head +
			ok + " := " + errV + " == nil\n" +
			"var " + actual + " any = " + v + "\n" +
			"if " + errV + " != nil {\n" +
			v + ", " + actual + " = nil, " + errV + "\n" +
			"}\n"
		return valCode{Setup: setup, V: v, OK: ok, Actual: actual}
	case needOK:
		setup := head + ok + " := " + errV + " == nil\n"
		return valCode{Setup: setup, V: v, OK: ok}
	default:
		setup := head +
			"if " + errV + " != nil {\n" +
			v + " = nil\n" +
			"}\n"
		return valCode{Setup: setup, V: v}
	}
}
