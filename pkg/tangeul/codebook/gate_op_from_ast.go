//ff:func feature=tangl type=util control=selection
//ff:what GateOpFromAST — map an ast.Certainty.Op string to its word.GateOp
package codebook

import (
	"fmt"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/word"
)

// GateOpFromAST maps a canonical ast.Certainty.Op string ("at least", "above",
// "less than", "at most") to its word.GateOp (§3.2). An unrecognized string is
// a hard error — never a silent fallback.
func GateOpFromAST(op string) (word.GateOp, error) {
	switch op {
	case "at least":
		return word.GateAtLeast, nil
	case "above":
		return word.GateAbove, nil
	case "less than":
		return word.GateLessThan, nil
	case "at most":
		return word.GateAtMost, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Certainty op %q", op)
	}
}
