//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what wrapClosures — wraps each compiled rule as a public (ctx, specs) rule func
package eval

import (
	"fmt"
	"reflect"

	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
	"github.com/park-jun-woo/toulmin/go/pkg/toulmin"
)

// wrapClosures wraps every compiled rule closure with a fresh-scope seeder,
// producing the public rule-func map keyed by rule name.
func wrapClosures(rules []ast.InlineRule, compiled map[string]evalFunc) map[string]func(toulmin.Context, toulmin.Specs) (bool, any) {
	return wrapClosuresWithSpecs(rules, compiled, nil)
}

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

func mergeRuleSpecs(owned, supplied toulmin.Specs) (toulmin.Specs, error) {
	merged := append(toulmin.Specs(nil), owned...)
	for _, candidate := range supplied {
		existing := merged.Find(candidate.SpecName())
		if existing == nil {
			merged = append(merged, candidate)
			continue
		}
		if !reflect.DeepEqual(existing, candidate) {
			return nil, fmt.Errorf("conflicting Spec %q supplied to rule", candidate.SpecName())
		}
	}
	return merged, nil
}
