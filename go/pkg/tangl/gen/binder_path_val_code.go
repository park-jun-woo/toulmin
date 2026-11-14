//ff:func feature=tangl type=codegen control=sequence
//ff:what binderPathValCode — emits the code resolving a binder-headed Path operand
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// binderPathValCode emits the code resolving a binder-headed Path: the
// binder-only path is the bound element variable itself (always ok), and a
// multi-step path walks the element via tanglWalk — the binder branch of
// pkg/tangl/eval.pathVal.
func binderPathValCode(p ast.Path, b binderInfo, prefix string, needOK bool) valCode {
	if len(p.Steps) == 1 {
		boundOK := ""
		if needOK {
			boundOK = "true"
		}
		return valCode{V: b.Elem, OK: boundOK, Actual: b.Elem}
	}
	v := prefix + "v"
	ok := prefix + "ok"
	okLHS := ok
	if !needOK {
		okLHS = "_"
		ok = ""
	}
	setup := v + ", " + okLHS + " := tanglWalk(" + b.Elem + ", " + quoteSteps(p.Steps[1:]) + ")\n"
	return valCode{Setup: setup, V: v, OK: ok, Actual: v}
}
