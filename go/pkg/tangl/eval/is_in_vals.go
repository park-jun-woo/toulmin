//ff:func feature=tangl type=engine control=sequence
//ff:what isInVals — "is in" membership of a resolved lhs in a structured string list or a comma-list rhs
package eval

import "fmt"

// isInVals reports whether a resolved lhs is a member of rhs. A missing lhs is
// never a member. A structured rhs exposing StringList() is matched exactly
// against the lhs "%v" form; any other rhs falls back to compareIn's comma-list
// membership over both "%v" forms.
func isInVals(lhs any, lhsOK bool, rhs any) bool {
	if !lhsOK {
		return false
	}
	if list, ok := rhs.(interface{ StringList() []string }); ok {
		return containsString(list.StringList(), fmt.Sprintf("%v", lhs))
	}
	return compareIn(fmt.Sprintf("%v", lhs), fmt.Sprintf("%v", rhs))
}
