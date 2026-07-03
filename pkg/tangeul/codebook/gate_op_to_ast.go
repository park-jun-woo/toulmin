//ff:func feature=tangl type=util control=selection
//ff:what GateOpToAST — map a word.GateOp back to its ast.Certainty.Op string
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// GateOpToAST maps a word.GateOp back to its canonical ast.Certainty.Op string
// (§3.2), the inverse of GateOpFromAST. An unrecognized op is a hard error —
// never a silent fallback.
func GateOpToAST(op word.GateOp) (string, error) {
	switch op {
	case word.GateAtLeast:
		return "at least", nil
	case word.GateAbove:
		return "above", nil
	case word.GateLessThan:
		return "less than", nil
	case word.GateAtMost:
		return "at most", nil
	default:
		return "", fmt.Errorf("tangeul/codebook: unknown GateOp %d", uint8(op))
	}
}
