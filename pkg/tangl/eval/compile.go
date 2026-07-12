//ff:func feature=tangl type=engine control=sequence
//ff:what Compile — compiles one InlineRule condition into a rule-func closure
package eval

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// Compile compiles a single InlineRule's condition tree into a closure with
// the toulmin.RuleMeta.Fn signature. Every Call is bound to res at compile
// time; an unresolved reference is returned as an error and never becomes a
// runtime lookup.
func Compile(r ast.InlineRule, res Resolver) (func(toulmin.Context, toulmin.Specs) (bool, any), error) {
	ef, err := compileExpr(r.Cond, res)
	if err != nil {
		return nil, err
	}
	return func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
		return ef(newScope(ctx, specs))
	}, nil
}
