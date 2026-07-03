//ff:func feature=tangl type=util control=selection
//ff:what LogicOpFromAST — map an ast.Logic.Op string ("and"/"or") to its LogicOp code
package codebook

import "fmt"

// LogicOpFromAST maps a canonical ast.Logic.Op string ("and" or "or") to its
// LogicOp code. An unrecognized string is a hard error — never a silent
// fallback.
func LogicOpFromAST(op string) (LogicOp, error) {
	switch op {
	case "and":
		return LogicAnd, nil
	case "or":
		return LogicOr, nil
	default:
		return 0, fmt.Errorf("tangeul/codebook: unknown Logic op %q", op)
	}
}
