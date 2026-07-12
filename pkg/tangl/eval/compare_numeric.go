//ff:func feature=tangl type=engine control=sequence
//ff:what compareNumeric — parses both operands as floats and applies a numeric operator
package eval

import (
	"fmt"
	"strconv"
)

// compareNumeric parses lhs and rhs from their "%v" string form as float64 and
// applies the numeric operator. If either operand does not parse, it returns
// false.
func compareNumeric(op string, lhs, rhs any) bool {
	fv, ferr := strconv.ParseFloat(fmt.Sprintf("%v", lhs), 64)
	tv, terr := strconv.ParseFloat(fmt.Sprintf("%v", rhs), 64)
	if ferr != nil || terr != nil {
		return false
	}
	return numericOp(op, fv, tv)
}
