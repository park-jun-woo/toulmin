//ff:func feature=tangl type=util control=sequence
//ff:what compensationStackOf — retrieves the compensation stack from ctx, nil if absent or malformed
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// compensationStackOf returns the *compensationStack stored in ctx, or nil if
// InitCompensation was never called (or the stack was cleared/miskeyed).
// It never panics, even on a ctx that has never seen the compensation key.
func compensationStackOf(ctx toulmin.Context) *compensationStack {
	v, ok := ctx.Get(compensationKey)
	if !ok {
		return nil
	}
	st, ok := v.(*compensationStack)
	if !ok {
		return nil
	}
	return st
}
