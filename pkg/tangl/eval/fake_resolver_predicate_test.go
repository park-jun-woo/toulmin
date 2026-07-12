//ff:func feature=tangl type=engine control=sequence
//ff:what Predicate — test-Resolver predicate lookup keyed by alias and name
package eval

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Predicate looks up a predicate function by alias and name.
func (f fakeResolver) Predicate(alias, name string) (func(toulmin.Context, toulmin.Specs) (bool, any), bool) {
	fn, ok := f.predicates[alias+"\x00"+name]
	return fn, ok
}
