//ff:func feature=tangl type=util control=selection
//ff:what GateOpFromAST — map an ast.Certainty.Op string to its toulmin.GateOp
package codebook

import (
	"fmt"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

// GateOpFromAST maps a canonical ast.Certainty.Op string ("at least", "above",
// "less than", "at most") to its toulmin.GateOp (§3.2). An unrecognized string is
// a hard error — never a silent fallback.
func GateOpFromAST(op string) (toulmin.GateOp, error) {
	switch op {
	case "at least":
		return toulmin.GateAtLeast, nil
	case "above":
		return toulmin.GateAbove, nil
	case "less than":
		return toulmin.GateLessThan, nil
	case "at most":
		return toulmin.GateAtMost, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Certainty op %q", op)
	}
}
