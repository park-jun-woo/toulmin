//ff:func feature=tangl type=engine control=sequence
//ff:what compileEither — compiles an Either group into a first-true closure
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// compileEither compiles an Either group. Its short-circuit semantics are
// identical to Logic "or" (first true wins, last witness on all-false), so it
// reuses orClosure.
func compileEither(e ast.Either, res Resolver) (evalFunc, error) {
	terms, err := compileTerms(e.Terms, res)
	if err != nil {
		return nil, err
	}
	return orClosure(terms), nil
}
