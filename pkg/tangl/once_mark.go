//ff:func feature=tangl type=engine control=sequence
//ff:what OnceMark — marks a once guard as fired, consumed only after a successful action
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// OnceMark marks the once guard identified by key as fired. Callers (codegen'd
// RunOn handlers) must call this only after the guarded action has succeeded,
// so a failed or certainty-gated attempt is retried on the next tick.
func OnceMark(ctx toulmin.Context, key string) {
	ctx.Set(key, true)
}
