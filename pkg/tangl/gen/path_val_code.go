//ff:func feature=tangl type=codegen control=sequence
//ff:what pathValCode — emits the code resolving a Path operand to (value, ok)
package gen

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// pathValCode emits the code resolving a Path operand, mirroring
// pkg/tangl/eval.pathVal: a path whose first step names an in-scope exists
// binder resolves from the bound element (binderPathValCode), any other path
// is ctx-rooted via tanglCtxVal. The prefix keeps emitted variable names
// distinct per operand position; ok is emitted as a blank assign when the
// caller does not need it.
func pathValCode(p ast.Path, env map[string]binderInfo, prefix string, needOK bool) valCode {
	if len(p.Steps) > 0 {
		if b, bound := env[p.Steps[0]]; bound {
			return binderPathValCode(p, b, prefix, needOK)
		}
	}
	v := prefix + "v"
	ok := prefix + "ok"
	okLHS := ok
	if !needOK {
		okLHS = "_"
		ok = ""
	}
	setup := v + ", " + okLHS + " := tanglCtxVal(ctx, " + quoteSteps(p.Steps) + ")\n"
	return valCode{Setup: setup, V: v, OK: ok, Actual: v}
}
