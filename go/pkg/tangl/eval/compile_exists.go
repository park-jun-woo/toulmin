//ff:func feature=tangl type=engine control=sequence
//ff:what compileExists — compiles an Exists quantifier into a first-satisfying closure
package eval

import (
	"reflect"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// compileExists compiles an Exists quantifier. At runtime it resolves Source
// to a slice/array, iterates elements in index order binding each to Binder,
// and returns true with the satisfying element's witness at the first element
// whose Cond holds. A no-value or non-slice source yields zero elements; when
// no element satisfies it returns a Fact whose Actual is the number of
// elements examined.
func compileExists(e ast.Exists, res Resolver) (evalFunc, error) {
	src, err := compileOperandValue(e.Source, res)
	if err != nil {
		return nil, err
	}
	cond, err := compileExpr(e.Cond, res)
	if err != nil {
		return nil, err
	}
	source := e.Source
	binder := e.Binder
	return func(sc *scope) (bool, any) {
		expected := "any of " + reprOperand(source, sc) + " as " + binder + " satisfies"
		sv, ok, serr := src(sc)
		if !ok || serr != nil {
			return false, toulmin.Fact{Where: reprOperand(source, sc), Expected: expected, Actual: 0}
		}
		rv := reflect.ValueOf(sv)
		if k := rv.Kind(); k != reflect.Slice && k != reflect.Array {
			return false, toulmin.Fact{Where: reprOperand(source, sc), Expected: expected, Actual: 0}
		}
		prev, had := sc.binders[binder]
		n := rv.Len()
		for i := 0; i < n; i++ {
			sc.binders[binder] = &binding{source: source, index: i, value: rv.Index(i).Interface()}
			b, w := cond(sc)
			if b {
				if had {
					sc.binders[binder] = prev
				} else {
					delete(sc.binders, binder)
				}
				return true, w
			}
		}
		if had {
			sc.binders[binder] = prev
		} else {
			delete(sc.binders, binder)
		}
		return false, toulmin.Fact{Where: reprOperand(source, sc), Expected: expected, Actual: n}
	}, nil
}
