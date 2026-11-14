//ff:func feature=tangl type=codegen control=sequence
//ff:what equivResolver.Predicate — resolves the bare isSpam host predicate
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// Predicate resolves the cross-test document's bare isSpam predicate call.
func (equivResolver) Predicate(alias, name string) (func(toulmin.Context, toulmin.Specs) (bool, any), bool) {
	if alias == "" && name == "isSpam" {
		return equivIsSpam, true
	}
	return nil, false
}
