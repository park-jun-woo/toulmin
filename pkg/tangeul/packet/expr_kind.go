//ff:type feature=tangl type=model
//ff:what ExprKind — the 3-bit record kind of a flattened Rules condition node (§design 1)
package packet

// ExprKind is the 3-bit kind of an EXPR record. It is packed across a header
// word as kind = (bit13 << 2) | bits15-14, leaving the low 13 bits (bits 12-0)
// for the child count (max 8191). The original four kinds (0-3) all have bit13
// unset, so their header bit patterns are byte-identical to the earlier 2-bit
// layout (kind in bits 15-14, count in bits 13-0) as long as the child count
// stays below 8192 — which the encoder now guards. The prefix-order record
// stream plus per-record child counts uniquely rebuild the ast.Expr tree.
type ExprKind uint16

const (
	// ExprCompare is a leaf `<field> <op> <value>` comparison (ast.Compare)
	// with no operands (Lhs == nil && Rhs == nil).
	ExprCompare ExprKind = iota
	// ExprLogic is an "and"/"or" combiner (ast.Logic).
	ExprLogic
	// ExprNot is a negation (ast.Not), always one child.
	ExprNot
	// ExprEither is an `either` group (ast.Either).
	ExprEither
	// ExprExists is an "any of ... satisfies" quantifier (ast.Exists), one child
	// (the Cond subtree). bit13=1, bits15-14=00.
	ExprExists
	// ExprCall is a host-function call in predicate position (ast.Call as an
	// Expr). bit13=1, bits15-14=01.
	ExprCall
	// ExprCompareX is a comparison carrying an operand in Lhs and/or Rhs
	// (ast.Compare with Lhs != nil || Rhs != nil). bit13=1, bits15-14=10.
	ExprCompareX
	// exprKindReserved (7) is reserved for a future record kind. bit13=1,
	// bits15-14=11.
	exprKindReserved
)
