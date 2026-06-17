//ff:func feature=engine type=engine control=sequence
//ff:what safeCallHandler — executes a node handler with panic recovery
package toulmin

import "fmt"

// safeCallHandler runs h with panic recovery. Returns error if h panics or returns an error.
func safeCallHandler(h NodeHandler, ctx Context, ev NodeEvent, view RunView) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handler panicked: %v", r)
		}
	}()
	return h(ctx, ev, view)
}
