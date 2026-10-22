//ff:func feature=tangl type=util control=selection
//ff:what CompareOpToAST — map a CompareOp code back to its ast.Compare.Op string
package codebook

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// CompareOpToAST maps a CompareOp code back to its canonical ast.Compare.Op
// string, the inverse of CompareOpFromAST. An unrecognized code is a hard error
// — never a silent fallback.
func CompareOpToAST(op toulmin.CompareOp) (string, error) {
	switch op {
	case toulmin.CompareIsEmpty:
		return "is empty", nil
	case toulmin.CompareIsNotEmpty:
		return "is not empty", nil
	case toulmin.CompareEquals:
		return "equals", nil
	case toulmin.CompareIsIn:
		return "is in", nil
	case toulmin.CompareIsGreaterThan:
		return "is greater than", nil
	case toulmin.CompareIsLessThan:
		return "is less than", nil
	case toulmin.CompareIsAtMost:
		return "is at most", nil
	case toulmin.CompareIsAtLeast:
		return "is at least", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown toulmin.CompareOp 0x%04X", uint16(op))
	}
}
