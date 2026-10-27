//ff:func feature=tangl type=util control=selection
//ff:what LogicOpToAST — map a LogicOp code back to its ast.Logic.Op string
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// LogicOpToAST maps a LogicOp code back to its canonical ast.Logic.Op string,
// the inverse of LogicOpFromAST. An unrecognized code is a hard error — never a
// silent fallback.
func LogicOpToAST(op toulmin.LogicOp) (string, error) {
	switch op {
	case toulmin.LogicAnd:
		return "and", nil
	case toulmin.LogicOr:
		return "or", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown toulmin.LogicOp 0x%04X", uint16(op))
	}
}
