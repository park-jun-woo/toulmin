//ff:type feature=tangl type=model
//ff:what ExprKind — the 2-bit record kind of a flattened Rules condition node (§design 4)
package packet

// ExprKind is the 2-bit kind stored in bits 15-14 of an EXPR record header word
// (§design 4). The low 14 bits of that header hold the child count, so the
// prefix-order record stream plus per-record child counts uniquely rebuild the
// ast.Expr tree.
type ExprKind uint16

const (
	// ExprCompare is a leaf `<field> <op> <value>` comparison (ast.Compare).
	ExprCompare ExprKind = iota
	// ExprLogic is an "and"/"or" combiner (ast.Logic).
	ExprLogic
	// ExprNot is a negation (ast.Not), always one child.
	ExprNot
	// ExprEither is an `either` group (ast.Either).
	ExprEither
)
