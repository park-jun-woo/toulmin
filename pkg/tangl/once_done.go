//ff:func feature=tangl type=engine control=sequence
//ff:what OnceDone — reports whether a once-guarded do has already fired successfully
package tangl

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

// OnceDone reports whether the once guard identified by key has already been
// marked (via OnceMark) on ctx. Keys follow the codegen convention
// "once:<subject>.<case>.<node>#<n>".
func OnceDone(ctx toulmin.Context, key string) bool {
	v, ok := ctx.Get(key)
	if !ok {
		return false
	}
	done, ok := v.(bool)
	if !ok {
		return false
	}
	return done
}
