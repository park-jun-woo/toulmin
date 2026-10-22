//ff:func feature=tangl type=engine control=sequence
//ff:what newScope — seeds a fresh evaluation scope from a ctx and specs
package eval

import "github.com/park-jun-woo/toulmin/go/pkg/toulmin"

// newScope creates a fresh scope with an empty binder set for one top-level
// rule evaluation.
func newScope(ctx toulmin.Context, specs toulmin.Specs) *scope {
	return &scope{ctx: ctx, specs: specs, binders: map[string]*binding{}}
}
