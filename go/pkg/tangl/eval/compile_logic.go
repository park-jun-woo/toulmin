//ff:func feature=tangl type=engine control=sequence
//ff:what compileLogic — compiles an and/or Logic node into a short-circuit closure
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"

// compileLogic compiles a Logic node: it compiles the terms and selects the
// short-circuit closure for its operator.
func compileLogic(l ast.Logic, res Resolver) (evalFunc, error) {
	terms, err := compileTerms(l.Terms, res)
	if err != nil {
		return nil, err
	}
	return logicClosure(l.Op, terms)
}
