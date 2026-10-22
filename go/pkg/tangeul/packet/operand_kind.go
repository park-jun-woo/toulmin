//ff:type feature=tangl type=model
//ff:what OperandKind — how a value-producing operand is stored in an EXPR record (§design 2)
package packet

// OperandKind discriminates an ast.Operand slot in an EXPR record (§design 2).
// It is an EXPR-payload record format, not a codebook vocabulary code, so it
// lives in packet next to ExprKind and ValueKind. All operand strings are
// interned into the packet's own label table (empty strings included), so no
// sentinel index is needed for an empty alias.
type OperandKind uint16

const (
	// OperandPath is a ctx data path (ast.Path): [0][step count][label idx × step].
	OperandPath OperandKind = iota
	// OperandSpec is a with-injected spec reference (ast.SpecRef):
	// [1][step count][label idx × step].
	OperandSpec
	// OperandCall is a host-function call in value position (ast.Call):
	// [2][alias label idx][name label idx]; a local call interns "" as its alias.
	OperandCall
)
