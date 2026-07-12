//ff:func feature=tangl type=codegen control=sequence
//ff:what specValCode — emits the code resolving a SpecRef operand to (value, ok)
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// specValCode emits the code resolving a SpecRef operand via the tanglSpecVal
// helper (Specs.Find on the first step, then a tanglWalk of the rest),
// mirroring pkg/tangl/eval.specVal. The prefix keeps emitted variable names
// distinct per operand position; ok is emitted as a blank assign when the
// caller does not need it.
func specValCode(sr ast.SpecRef, prefix string, needOK bool) valCode {
	v := prefix + "v"
	ok := prefix + "ok"
	okLHS := ok
	if !needOK {
		okLHS = "_"
		ok = ""
	}
	setup := v + ", " + okLHS + " := tanglSpecVal(specs, " + quoteSteps(sr.Steps) + ")\n"
	return valCode{Setup: setup, V: v, OK: ok, Actual: v}
}
