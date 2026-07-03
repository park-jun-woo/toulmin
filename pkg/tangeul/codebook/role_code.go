//ff:type feature=tangl type=model
//ff:what RoleCode — the 2-bit RuleNode role subcode (codebook §2.2)
package codebook

// RoleCode is the 2-bit RuleNode role subcode stored in bits 15-14 of the Node
// header word (TANGEUL §2.2, §5.1). It maps 1:1 to ast.Role. Free-area L1
// codebook value, FROZEN append-only.
type RoleCode uint16

// The three TANGEUL role subcodes (§2.2).
const (
	// RoleGeneral is a general rule — engine Rule (Warrant), asserts (0).
	RoleGeneral RoleCode = 0
	// RoleCounter is a counter rule — engine Counter (Rebuttal), reverses (1).
	RoleCounter RoleCode = 1
	// RoleExcept is an except rule — engine Except (Defeater), attenuates (2).
	RoleExcept RoleCode = 2
)
