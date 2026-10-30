//ff:func feature=tangl type=engine control=sequence
//ff:what wrapClosures — wraps each compiled rule as a public (ctx, specs) rule func
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// wrapClosures wraps every compiled rule closure with a fresh-scope seeder,
// producing the public rule-func map keyed by rule name.
func wrapClosures(rules []ast.InlineRule, compiled map[string]evalFunc) map[string]func(toulmin.Context, toulmin.Specs) (bool, any) {
	return wrapClosuresWithSpecs(rules, compiled, nil)
}
