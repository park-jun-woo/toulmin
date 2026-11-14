//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what wrapClosuresWithSpecs — wrap each compiled rule as a (ctx, specs) func that merges its rule-owned Specs
package eval

import (
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// wrapClosuresWithSpecs wraps every compiled rule closure with a fresh-scope
// seeder. The rule's derived (rule-owned) Specs are merged with the Specs
// supplied at call time; a merge conflict makes the rule return (false, err).
func wrapClosuresWithSpecs(rules []ast.InlineRule, compiled map[string]evalFunc, derived map[string]toulmin.Specs) map[string]func(toulmin.Context, toulmin.Specs) (bool, any) {
	out := map[string]func(toulmin.Context, toulmin.Specs) (bool, any){}
	for _, r := range rules {
		ef := compiled[r.Name]
		owned := derived[r.Name]
		out[r.Name] = func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
			merged, err := mergeRuleSpecs(owned, specs)
			if err != nil {
				return false, err
			}
			return ef(newScope(ctx, merged))
		}
	}
	return out
}
