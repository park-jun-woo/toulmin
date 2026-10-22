//ff:func feature=tangl type=util control=selection
//ff:what QuantOpFromAST — map a Rules quantifier string ("any") to its QuantOp code
package codebook

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// QuantOpFromAST maps a canonical Rules quantifier string ("any", the sole
// implemented quantifier of an ast.Exists node) to its QuantOp code. An
// unrecognized string is a hard error — never a silent fallback.
func QuantOpFromAST(op string) (toulmin.QuantOp, error) {
	switch op {
	case "any":
		return toulmin.QuantAny, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Quant op %q", op)
	}
}
