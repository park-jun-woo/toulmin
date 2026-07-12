//ff:func feature=tangl type=util control=selection
//ff:what QuantOpToAST — map a QuantOp code back to its Rules quantifier string
package codebook

import "fmt"

// QuantOpToAST maps a QuantOp code back to its canonical Rules quantifier
// string, the inverse of QuantOpFromAST. An unrecognized code is a hard error —
// never a silent fallback.
func QuantOpToAST(op QuantOp) (string, error) {
	switch op {
	case QuantAny:
		return "any", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown QuantOp 0x%04X", uint16(op))
	}
}
