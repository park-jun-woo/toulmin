//ff:func feature=tangl type=codegen control=sequence
//ff:what equivResolver.Value — resolves the bare nlpScore host value function
package gen

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// Value resolves the cross-test document's bare nlpScore value call.
func (equivResolver) Value(alias, name string) (func(toulmin.Context, toulmin.Specs) (any, error), bool) {
	if alias == "" && name == "nlpScore" {
		return equivNlpScore, true
	}
	return nil, false
}
