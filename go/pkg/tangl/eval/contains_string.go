//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what containsString — exact membership of a string in a list
package eval

// containsString reports whether want equals one of items exactly.
func containsString(items []string, want string) bool {
	for _, item := range items {
		if item == want {
			return true
		}
	}
	return false
}
