//ff:func feature=tangl type=util control=selection
//ff:what LogicOpToAST — map a LogicOp code back to its ast.Logic.Op string
package codebook

import "fmt"

// LogicOpToAST maps a LogicOp code back to its canonical ast.Logic.Op string,
// the inverse of LogicOpFromAST. An unrecognized code is a hard error — never a
// silent fallback.
func LogicOpToAST(op LogicOp) (string, error) {
	switch op {
	case LogicAnd:
		return "and", nil
	case LogicOr:
		return "or", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown LogicOp 0x%04X", uint16(op))
	}
}
