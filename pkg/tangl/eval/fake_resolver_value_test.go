//ff:func feature=tangl type=engine control=sequence
//ff:what Value — test-Resolver value lookup keyed by alias and name
package eval

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Value looks up a value function by alias and name.
func (f fakeResolver) Value(alias, name string) (func(toulmin.Context, toulmin.Specs) (any, error), bool) {
	fn, ok := f.values[alias+"\x00"+name]
	return fn, ok
}
