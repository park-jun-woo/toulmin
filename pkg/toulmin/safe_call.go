//ff:func feature=engine type=engine control=sequence
//ff:what safeCall — executes rule function with panic recovery
package toulmin

import "fmt"

// safeCall runs fn with panic recovery. Returns error if fn panics.
func safeCall(fn func(any, any, Backing) (bool, any), claim, ground any, backing Backing) (activated bool, evidence any, err error) {
	defer func() {
		if r := recover(); r != nil {
			activated = false
			evidence = nil
			err = fmt.Errorf("rule panicked: %v", r)
		}
	}()
	activated, evidence = fn(claim, ground, backing)
	return
}
