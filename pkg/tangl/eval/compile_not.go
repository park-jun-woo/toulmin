//ff:func feature=tangl type=engine control=sequence
//ff:what compileNot — compiles a Not node into a verdict-inverting closure
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// compileNot compiles a Not node: it evaluates the inner expression and
// inverts the verdict, propagating the inner witness unchanged.
func compileNot(n ast.Not, res Resolver) (evalFunc, error) {
	inner, err := compileExpr(n.Term, res)
	if err != nil {
		return nil, err
	}
	return func(sc *scope) (bool, any) {
		b, w := inner(sc)
		return !b, w
	}, nil
}
