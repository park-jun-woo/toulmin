//ff:func feature=approve type=util control=selection
//ff:what equalAny — compares two values by identity for maps and by value otherwise
package approve

// equalAny compares two values by identity for reference types (maps,
// interfaces, pointers) and by value otherwise, avoiding a reflect.DeepEqual
// dependency for this straightforward field-copy check.
func equalAny(a, b any) bool {
	switch av := a.(type) {
	case map[string]any:
		bv, ok := b.(map[string]any)
		if !ok || len(av) != len(bv) {
			return false
		}
		return mapsEqual(av, bv)
	default:
		return a == b
	}
}
