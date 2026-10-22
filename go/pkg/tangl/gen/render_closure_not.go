//ff:func feature=tangl type=codegen control=sequence pattern=recursive
//ff:what renderClosureNot — emits a Not node as a tanglNot closure literal
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// renderClosureNot emits a Not node as a func() (bool, any) closure literal
// wrapping the inner term in tanglNot, which inverts the verdict and
// propagates the inner witness unchanged — matching pkg/tangl/eval.compileNot.
func renderClosureNot(n ast.Not, env map[string]binderInfo, depth int) (string, error) {
	inner, err := renderClosureExpr(n.Term, env, depth)
	if err != nil {
		return "", err
	}
	return "func() (bool, any) {\nreturn tanglNot(" + inner + ")\n}", nil
}
