//ff:func feature=tangl type=engine control=sequence
//ff:what CommitCompensation — discards the current pass's compensation stack on success
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// CommitCompensation discards the armed compensation stack after a successful Run
// pass, so a later Compensate call on the same ctx becomes a no-op.
func CommitCompensation(ctx toulmin.Context) {
	ctx.Set(compensationKey, nil)
}
