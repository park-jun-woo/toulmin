//ff:func feature=tangl type=engine control=sequence
//ff:what compileCall — compiles a predicate-position Call into a holds closure
package eval

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// compileCall compiles a predicate-position Call (`x.y holds`). It binds the
// call to res.Predicate at compile time; an unresolved reference is an error.
// At runtime, if the host returns non-nil evidence that evidence is the
// witness (passed through unchanged); otherwise the witness is a Fact recording
// the call and its boolean result.
func compileCall(c ast.Call, res Resolver) (evalFunc, error) {
	fn, ok := res.Predicate(c.Fn.Alias, c.Fn.Name)
	if !ok {
		return nil, fmt.Errorf("eval: unresolved predicate %s", callRepr(c))
	}
	repr := callRepr(c)
	return func(sc *scope) (bool, any) {
		b, ev := fn(sc.ctx, sc.specs)
		if ev != nil {
			return b, ev
		}
		return b, toulmin.Fact{Where: repr, Expected: "holds", Actual: b}
	}, nil
}
