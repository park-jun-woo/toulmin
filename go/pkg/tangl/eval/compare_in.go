//ff:func feature=tangl type=engine control=iteration dimension=1
//ff:what compareIn — reports whether lhs equals any comma-separated item of rhs
package eval

import "strings"

// compareIn reports whether the lhs string equals any comma-separated,
// space-trimmed item of the rhs string.
func compareIn(lhs, rhs string) bool {
	for _, item := range strings.Split(rhs, ",") {
		if strings.TrimSpace(item) == lhs {
			return true
		}
	}
	return false
}
