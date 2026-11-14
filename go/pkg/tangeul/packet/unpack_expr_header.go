//ff:func feature=tangl type=util control=sequence
//ff:what unpackExprHeader — split an EXPR record header word into its 3-bit ExprKind and 13-bit child count
package packet

import "github.com/geul-org/gogeul/pkg/core"

// unpackExprHeader is the inverse of packExprHeader: it reassembles the 3-bit
// ExprKind from bits 15-14 (low two bits) and bit 13 (high bit), and reads the
// child count from bits 12-0 (§design 1). Legacy headers (kinds 0-3, count <
// 8192) decode identically because their bit 13 was always zero.
func unpackExprHeader(header core.Word) (ExprKind, int) {
	h := uint16(header)
	kind := ExprKind((h>>13&1)<<2 | h>>14)
	count := int(h & exprCountMax)
	return kind, count
}
