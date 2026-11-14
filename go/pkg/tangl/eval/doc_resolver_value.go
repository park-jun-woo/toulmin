//ff:func feature=tangl type=engine control=sequence
//ff:what Value — delegates value-Call resolution straight to the host Resolver
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// Value delegates value-position Call resolution straight to the host Resolver:
// document rules are predicates, not value functions.
func (d docResolver) Value(alias, name string) (func(toulmin.Context, toulmin.Specs) (any, error), bool) {
	return d.host.Value(alias, name)
}
