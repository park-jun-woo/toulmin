//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestUnpackExprHeader — golden header words split into every ExprKind with legacy bit-pattern identity
package packet

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/core"
)

// TestUnpackExprHeader locks the EXPR header decode (§design 1) against golden
// words: the four legacy kinds read their old `kind<<14 | count` pattern (bit
// 13 zero), the bit-13 kinds reassemble as (bit13 << 2) | bits15-14, and the
// child count comes from bits 12-0 up to its 8191 maximum.
func TestUnpackExprHeader(t *testing.T) {
	cases := []struct {
		header uint16
		kind   ExprKind
		count  int
	}{
		{0x0000, ExprCompare, 0},
		{0x4003, ExprLogic, 3},
		{0x8001, ExprNot, 1},
		{0xC002, ExprEither, 2},
		{0x2001, ExprExists, 1},
		{0x6000, ExprCall, 0},
		{0xA000, ExprCompareX, 0},
		{0xE000, exprKindReserved, 0},
		{0x5FFF, ExprLogic, exprCountMax},
	}
	for _, c := range cases {
		kind, count := unpackExprHeader(core.Word(c.header))
		if kind != c.kind || count != c.count {
			t.Fatalf("unpackExprHeader(0x%04X) = (%d, %d), want (%d, %d)", c.header, kind, count, c.kind, c.count)
		}
	}
}
