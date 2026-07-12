//ff:type feature=tangl type=model
//ff:what ValueKind — how a Compare value is stored in an EXPR record (NUM32 or Text)
package packet

// ValueKind discriminates a Compare value slot in an EXPR record (§design 4). A
// value that is an exact non-negative integer (its strconv round-trip is
// identical) is stored inline as a NUM32; a literal string value is stored as a
// Text reference into the packet's label table; and an ExprCompareX right side
// that is an ast.Operand (Rhs) is stored as an operand sub-record. This keeps
// the ast.Expr round-trip an exact identity for any value string or operand.
type ValueKind uint16

const (
	// ValueNum is an inline NUM32 numeric value.
	ValueNum ValueKind = iota
	// ValueText is a label-table Text reference value.
	ValueText
	// ValueOperand is an ast.Operand right side, encoded as an operand
	// sub-record (§design 3). Used only by ExprCompareX.
	ValueOperand
)
