//ff:func feature=tangl type=engine control=sequence
//ff:what PushCompensation — arms one compensation closure on the current pass's stack
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// PushCompensation arms fn to run during Compensate, called by a `do` handler
// immediately after its action succeeds. If InitCompensation was never called on
// ctx, a stack is created lazily so the push is never silently lost and no panic
// occurs. A nil fn is ignored.
func PushCompensation(ctx toulmin.Context, fn func() error) {
	if fn == nil {
		return
	}
	st := compensationStackOf(ctx)
	if st == nil {
		st = &compensationStack{}
		ctx.Set(compensationKey, st)
	}
	st.fns = append(st.fns, fn)
}
