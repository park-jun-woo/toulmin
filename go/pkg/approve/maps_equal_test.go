//ff:func feature=approve type=util control=iteration dimension=1
//ff:what mapsEqual — compares two string-keyed maps of any values for equality
package approve

// mapsEqual reports whether a and b contain the same set of keys mapped to
// equal values, comparing each value with ==. Callers must first verify
// len(a) == len(b) since this only checks that every key in a matches b.
func mapsEqual(a, b map[string]any) bool {
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
