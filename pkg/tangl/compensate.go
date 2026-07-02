//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what Compensate — runs the armed compensation stack in LIFO order, stopping at the first error
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// Compensate runs every armed compensation closure in reverse arming order (LIFO).
// It stops at the first error and returns it, leaving any remaining (deeper, earlier
// armed) closures unrun — per spec, promoting the failure to Review is the caller's
// responsibility. If ctx has no stack (InitCompensation never called, or the pass
// already committed), Compensate is a no-op and returns nil.
func Compensate(ctx toulmin.Context) error {
	st := compensationStackOf(ctx)
	if st == nil {
		return nil
	}
	for i := len(st.fns) - 1; i >= 0; i-- {
		if err := st.fns[i](); err != nil {
			return err
		}
	}
	return nil
}
