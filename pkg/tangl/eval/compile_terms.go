//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what compileTerms — compiles each Expr term into an evalFunc slice
package eval

import "github.com/park-jun-woo/toulmin/pkg/tangl/ast"

// compileTerms compiles every term expression into an evalFunc, returning the
// first compile error encountered.
func compileTerms(terms []ast.Expr, res Resolver) ([]evalFunc, error) {
	out := make([]evalFunc, len(terms))
	for i, t := range terms {
		ef, err := compileExpr(t, res)
		if err != nil {
			return nil, err
		}
		out[i] = ef
	}
	return out, nil
}
