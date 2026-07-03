//ff:func feature=tangl type=util control=selection
//ff:what CompareOpFromAST — map an ast.Compare.Op string to its CompareOp code
package codebook

import "fmt"

// CompareOpFromAST maps a canonical ast.Compare.Op string to its CompareOp code
// (Rules condition expressions). An unrecognized string is a hard error — never
// a silent fallback.
func CompareOpFromAST(op string) (CompareOp, error) {
	switch op {
	case "is empty":
		return CompareIsEmpty, nil
	case "is not empty":
		return CompareIsNotEmpty, nil
	case "equals":
		return CompareEquals, nil
	case "is in":
		return CompareIsIn, nil
	case "is greater than":
		return CompareIsGreaterThan, nil
	case "is less than":
		return CompareIsLessThan, nil
	case "is at most":
		return CompareIsAtMost, nil
	case "is at least":
		return CompareIsAtLeast, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Compare op %q", op)
	}
}
