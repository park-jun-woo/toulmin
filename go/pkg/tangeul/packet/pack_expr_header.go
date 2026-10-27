//ff:func feature=tangl type=util control=sequence
//ff:what packExprHeader — pack an EXPR record header word from a 3-bit ExprKind and a 13-bit child count
package packet

import (
	"fmt"

	"github.com/geul-org/gogeul/pkg/core"
)

// exprCountMax is the largest child count a 13-bit EXPR header field holds.
const exprCountMax = 0x1FFF

// packExprHeader packs an EXPR record header word: the 3-bit ExprKind is split
// as bits 15-14 (kind & 3) and bit 13 (kind >> 2), and the child count occupies
// bits 12-0 (§design 1). A count above 8191 (exprCountMax) would collide with
// bit 13 and silently corrupt the kind, so it is a hard error — this guard is
// what keeps the original four kinds' bit patterns byte-identical to the legacy
// layout for every stream the encoder produces.
func packExprHeader(kind ExprKind, count int) (core.Word, error) {
	if count < 0 || count > exprCountMax {
		return 0, fmt.Errorf("tangeul/packet: EXPR child count %d exceeds 13-bit max %d", count, exprCountMax)
	}
	hi := uint16(kind&3) << 14
	ext := uint16(kind>>2&1) << 13
	return core.Word(hi | ext | uint16(count)), nil
}
