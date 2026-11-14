//ff:func feature=engine type=engine control=sequence
//ff:what safeCall — executes rule function with panic recovery and verdict range check
package toulmin

import (
	"fmt"
	"math"
)

// safeCall runs fn with panic recovery. Returns error if fn panics or returns a
// verdict outside [-1, 1] (NaN included). On error the verdict is -1 (not applicable)
// and evidence is nil.
func safeCall(fn func(Context, Specs) (float64, any), ctx Context, specs Specs) (verdict float64, evidence any, err error) {
	defer func() {
		if r := recover(); r != nil {
			verdict = -1.0
			evidence = nil
			err = fmt.Errorf("rule panicked: %v", r)
		}
	}()
	verdict, evidence = fn(ctx, specs)
	if math.IsNaN(verdict) || verdict < -1.0 || verdict > 1.0 {
		return -1.0, nil, fmt.Errorf("rule verdict %v out of range [-1, 1]", verdict)
	}
	return verdict, evidence, nil
}
