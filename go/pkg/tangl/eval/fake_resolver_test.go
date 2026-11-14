//ff:type feature=tangl type=model
//ff:what fakeResolver — test Resolver backed by name-keyed predicate/value maps
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// fakeResolver is a test Resolver: it keys predicate and value functions by
// "alias\x00name" so both bare and alias-qualified lookups are exercised.
type fakeResolver struct {
	predicates map[string]func(toulmin.Context, toulmin.Specs) (bool, any)
	values     map[string]func(toulmin.Context, toulmin.Specs) (any, error)
}
