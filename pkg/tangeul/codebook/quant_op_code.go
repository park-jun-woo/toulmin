//ff:type feature=tangl type=model
//ff:what QuantOp — the Rules quantifier codes (ast.Exists), 0x03xx expression-operator block
package codebook

// QuantOp is the code for a collection quantifier used in Rules condition
// expressions. It maps to the canonical English quantifier of an ast.Exists
// node ("any"). Free-area L1 codebook value, FROZEN append-only.
//
// The code sits in the 0x03xx "expression operator" block, taking the next
// decimal slot (0x0321) after CompareOp 0x0301..0x0308 and LogicOp
// 0x0311..0x0312 — the same spacing rule the two neighbouring operator types
// use. The reserved quantifiers all/count/sum (grammar-reserved in Phase001,
// not yet implemented) are intentionally NOT assigned a code value here:
// append-only means a stamped-but-unimplemented code could never be retracted,
// so the phase that implements each will assign 0x0322 (all), 0x0323 (count)
// and 0x0324 (sum) at that time.
type QuantOp uint16

const (
	// QuantAny is the existential quantifier "any" (0x0321), ast.Exists.
	QuantAny QuantOp = 0x0321
)
