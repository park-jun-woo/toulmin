//ff:func feature=tangl type=engine control=selection
//ff:what compareVals — applies a tangl comparison operator to (lhs, rhs) values
package eval

import "fmt"

// compareVals applies one of the eight tangl comparison operators to a resolved
// lhs value (with its ok flag) against a rhs value. Both operands are compared
// in their fmt "%v" string form, generalizing the flat tanglCompare helper to
// value-vs-value comparison.
func compareVals(lhs any, lhsOK bool, op string, rhs any) bool {
	switch op {
	case "is empty":
		return !lhsOK || lhs == nil || fmt.Sprintf("%v", lhs) == ""
	case "is not empty":
		return lhsOK && lhs != nil && fmt.Sprintf("%v", lhs) != ""
	case "equals":
		return lhsOK && fmt.Sprintf("%v", lhs) == fmt.Sprintf("%v", rhs)
	case "is in":
		if !lhsOK {
			return false
		}
		return compareIn(fmt.Sprintf("%v", lhs), fmt.Sprintf("%v", rhs))
	case "is greater than", "is less than", "is at most", "is at least":
		if !lhsOK {
			return false
		}
		return compareNumeric(op, lhs, rhs)
	default:
		return false
	}
}
