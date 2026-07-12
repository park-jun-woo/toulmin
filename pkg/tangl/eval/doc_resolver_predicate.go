//ff:func feature=tangl type=engine control=selection
//ff:what Predicate — resolves a bare predicate to a sibling doc rule, else the host
package eval

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Predicate resolves a bare (alias-less) predicate Call to the compiled sibling
// document rule of that name when one exists, wrapping it to run in a fresh
// scope. Alias-qualified or non-rule references delegate to the host Resolver.
func (d docResolver) Predicate(alias, name string) (func(toulmin.Context, toulmin.Specs) (bool, any), bool) {
	if alias == "" && d.ruleSet[name] {
		nm := name
		return func(ctx toulmin.Context, specs toulmin.Specs) (bool, any) {
			return d.compiled[nm](newScope(ctx, specs))
		}, true
	}
	return d.host.Predicate(alias, name)
}
