//ff:func feature=tangl type=util control=selection
//ff:what LogicOpFromAST — map an ast.Logic.Op string ("and"/"or") to its LogicOp code
package codebook

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

// LogicOpFromAST maps a canonical ast.Logic.Op string ("and" or "or") to its
// LogicOp code. An unrecognized string is a hard error — never a silent
// fallback.
func LogicOpFromAST(op string) (toulmin.LogicOp, error) {
	switch op {
	case "and":
		return toulmin.LogicAnd, nil
	case "or":
		return toulmin.LogicOr, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Logic op %q", op)
	}
}
