//ff:func feature=tangl type=engine control=selection
//ff:what numericOp — compares two floats by the given tangl numeric operator
package eval

// numericOp compares fv against tv by the tangl numeric operator. The default
// case is "is at least" (>=), the only remaining numeric operator.
func numericOp(op string, fv, tv float64) bool {
	switch op {
	case "is greater than":
		return fv > tv
	case "is less than":
		return fv < tv
	case "is at most":
		return fv <= tv
	default:
		return fv >= tv
	}
}
