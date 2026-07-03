//ff:func feature=tangl type=util control=selection
//ff:what equalIgnoringLine — recursive reflect equality dispatch (skips Line, treats nil==empty slice)
package tangeul

import "reflect"

// equalIgnoringLine is the recursive worker behind EqualIgnoringLine: it walks two
// reflect.Values of identical type and reports structural equality. Pointers and
// interfaces recurse through their element; structs and slices delegate to their
// own helpers (which skip Line fields and treat zero-length slices as equal);
// every other kind (strings, numbers, bools) falls back to reflect.DeepEqual.
func equalIgnoringLine(x, y reflect.Value) bool {
	if x.Type() != y.Type() {
		return false
	}
	switch x.Kind() {
	case reflect.Pointer, reflect.Interface:
		if x.IsNil() || y.IsNil() {
			return x.IsNil() == y.IsNil()
		}
		return equalIgnoringLine(x.Elem(), y.Elem())
	case reflect.Struct:
		return equalStructIgnoringLine(x, y)
	case reflect.Slice:
		return equalSliceIgnoringLine(x, y)
	default:
		return reflect.DeepEqual(x.Interface(), y.Interface())
	}
}
