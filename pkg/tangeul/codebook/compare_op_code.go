//ff:type feature=tangl type=model
//ff:what CompareOp — the eight Rules comparison operator codes (ast.Compare.Op)
package codebook

// CompareOp is the code for a leaf comparison operator used in Rules condition
// expressions. It maps to the eight canonical English ast.Compare.Op strings.
// Free-area L1 codebook value, FROZEN append-only.
type CompareOp uint16

// The eight TANGEUL comparison operator codes, assigned 0x0301..0x0308.
const (
	// CompareIsEmpty is "is empty" (0x0301).
	CompareIsEmpty CompareOp = 0x0301
	// CompareIsNotEmpty is "is not empty" (0x0302).
	CompareIsNotEmpty CompareOp = 0x0302
	// CompareEquals is "equals" (0x0303).
	CompareEquals CompareOp = 0x0303
	// CompareIsIn is "is in" (0x0304).
	CompareIsIn CompareOp = 0x0304
	// CompareIsGreaterThan is "is greater than" (0x0305).
	CompareIsGreaterThan CompareOp = 0x0305
	// CompareIsLessThan is "is less than" (0x0306).
	CompareIsLessThan CompareOp = 0x0306
	// CompareIsAtMost is "is at most" (0x0307).
	CompareIsAtMost CompareOp = 0x0307
	// CompareIsAtLeast is "is at least" (0x0308).
	CompareIsAtLeast CompareOp = 0x0308
)
