//ff:func feature=engine type=engine control=sequence
//ff:what safeCall — executes rule function with panic recovery
package toulmin

import "fmt"

// safeCall runs fn with panic recovery. Returns error if fn panics.
func safeCall(fn func(Context, Specs) (bool, any), ctx Context, specs Specs) (activated bool, evidence any, err error) {
	defer func() {
		if r := recover(); r != nil {
			activated = false
			evidence = nil
			err = fmt.Errorf("rule panicked: %v", r)
		}
	}()
	activated, evidence = fn(ctx, specs)
	return
}
