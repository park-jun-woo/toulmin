//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what equalSliceIgnoringLine — element-wise slice equality (nil and empty compare equal)
package tangeul

import "reflect"

// equalSliceIgnoringLine compares two slices of identical type element by element,
// recursing through equalIgnoringLine. Length is the only nilness signal, so a nil
// and an empty slice compare equal (a representational, not semantic, difference).
func equalSliceIgnoringLine(x, y reflect.Value) bool {
	if x.Len() != y.Len() {
		return false
	}
	for i := 0; i < x.Len(); i++ {
		if !equalIgnoringLine(x.Index(i), y.Index(i)) {
			return false
		}
	}
	return true
}
