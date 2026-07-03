//ff:func feature=tangl type=util control=selection
//ff:what CompareOpToAST — map a CompareOp code back to its ast.Compare.Op string
package codebook

import "fmt"

// CompareOpToAST maps a CompareOp code back to its canonical ast.Compare.Op
// string, the inverse of CompareOpFromAST. An unrecognized code is a hard error
// — never a silent fallback.
func CompareOpToAST(op CompareOp) (string, error) {
	switch op {
	case CompareIsEmpty:
		return "is empty", nil
	case CompareIsNotEmpty:
		return "is not empty", nil
	case CompareEquals:
		return "equals", nil
	case CompareIsIn:
		return "is in", nil
	case CompareIsGreaterThan:
		return "is greater than", nil
	case CompareIsLessThan:
		return "is less than", nil
	case CompareIsAtMost:
		return "is at most", nil
	case CompareIsAtLeast:
		return "is at least", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown CompareOp 0x%04X", uint16(op))
	}
}
