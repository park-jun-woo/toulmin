//ff:type feature=tangl type=model
//ff:what LogicOp — the two Rules boolean combiner codes (ast.Logic.Op)
package codebook

// LogicOp is the code for a boolean combiner used in Rules condition
// expressions. It maps to the canonical English ast.Logic.Op strings "and" and
// "or". Free-area L1 codebook value, FROZEN append-only.
type LogicOp uint16

// The two TANGL boolean combiner codes, assigned 0x0311..0x0312.
const (
	// LogicAnd is "and" (0x0311).
	LogicAnd LogicOp = 0x0311
	// LogicOr is "or" (0x0312).
	LogicOr LogicOp = 0x0312
)
