//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what wrapClosures — wraps each compiled rule as a public (ctx, specs) rule func
package eval

import (
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/pkg/toulmin"
)

// wrapClosures wraps every compiled rule closure with a fresh-scope seeder,
// producing the public rule-func map keyed by rule name.
func wrapClosures(rules []ast.InlineRule, compiled map[string]evalFunc) map[string]func(toulmin.Context, toulmin.Specs) (bool, any) {
	out := map[string]func(toulmin.Context, toulmin.Specs) (bool, any){}
	for _, r := range rules {
		ef := compiled[r.Name]
		out[r.Name] = func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
			return ef(newScope(ctx, specs))
		}
	}
	return out
}
