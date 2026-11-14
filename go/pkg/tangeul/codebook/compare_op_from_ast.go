//ff:func feature=tangl type=util control=selection
//ff:what CompareOpFromAST — map an ast.Compare.Op string to its CompareOp code
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// CompareOpFromAST maps a canonical ast.Compare.Op string to its CompareOp code
// (Rules condition expressions). An unrecognized string is a hard error — never
// a silent fallback.
func CompareOpFromAST(op string) (toulmin.CompareOp, error) {
	switch op {
	case "is empty":
		return toulmin.CompareIsEmpty, nil
	case "is not empty":
		return toulmin.CompareIsNotEmpty, nil
	case "equals":
		return toulmin.CompareEquals, nil
	case "is in":
		return toulmin.CompareIsIn, nil
	case "is greater than":
		return toulmin.CompareIsGreaterThan, nil
	case "is less than":
		return toulmin.CompareIsLessThan, nil
	case "is at most":
		return toulmin.CompareIsAtMost, nil
	case "is at least":
		return toulmin.CompareIsAtLeast, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Compare op %q", op)
	}
}
