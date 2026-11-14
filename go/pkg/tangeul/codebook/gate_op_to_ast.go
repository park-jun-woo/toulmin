//ff:func feature=tangl type=util control=selection
//ff:what GateOpToAST — map a toulmin.GateOp back to its ast.Certainty.Op string
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// GateOpToAST maps a toulmin.GateOp back to its canonical ast.Certainty.Op string
// (§3.2), the inverse of GateOpFromAST. An unrecognized op is a hard error —
// never a silent fallback.
func GateOpToAST(op toulmin.GateOp) (string, error) {
	switch op {
	case toulmin.GateAtLeast:
		return "at least", nil
	case toulmin.GateAbove:
		return "above", nil
	case toulmin.GateLessThan:
		return "less than", nil
	case toulmin.GateAtMost:
		return "at most", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown toulmin.GateOp %d", uint8(op))
	}
}
