//ff:func feature=tangl type=util control=selection
//ff:what QuantOpToAST — map a QuantOp code back to its Rules quantifier string
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// QuantOpToAST maps a QuantOp code back to its canonical Rules quantifier
// string, the inverse of QuantOpFromAST. An unrecognized code is a hard error —
// never a silent fallback.
func QuantOpToAST(op toulmin.QuantOp) (string, error) {
	switch op {
	case toulmin.QuantAny:
		return "any", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown toulmin.QuantOp 0x%04X", uint16(op))
	}
}
