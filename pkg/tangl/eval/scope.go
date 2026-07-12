//ff:type feature=tangl type=model
//ff:what scope — runtime evaluation state (ctx, specs, in-scope exists binders)
package eval

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// scope carries the runtime state threaded through compiled closures: the
// evaluation Context, the rule's Specs, and the in-scope exists binder
// bindings keyed by binder name.
type scope struct {
	ctx     toulmin.Context
	specs   toulmin.Specs
	binders map[string]*binding
}
