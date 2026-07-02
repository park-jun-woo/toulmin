//ff:func feature=tangl type=engine control=sequence
//ff:what InitCompensation — arms a fresh, empty compensation stack for one Run pass
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// InitCompensation resets the compensation stack for a new top-level Run pass
// (Provides `run` wrapper per request, or Internal tick runner per tick).
// It must be called before the pass's do/undo handlers can push compensations.
func InitCompensation(ctx toulmin.Context) {
	ctx.Set(compensationKey, &compensationStack{})
}
