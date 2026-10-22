//ff:func feature=engine type=engine control=sequence
//ff:what TestFuncIDUniqueness — tests that different closures get distinct funcIDs
package toulmin

import (
	"testing"
)

func TestFuncIDUniqueness(t *testing.T) {
	fn1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	fn2 := func(ctx Context, specs Specs) (bool, any) { return false, nil }
	id1 := funcID(fn1)
	id2 := funcID(fn2)
	if id1 == id2 {
		t.Errorf("expected distinct funcIDs for different closures, both got %s", id1)
	}
}
