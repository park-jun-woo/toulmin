//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestPackExprHeader — header pack/unpack roundtrip, legacy bit-pattern identity, and the count guard
package packet

import (
	"testing"
)

// TestPackExprHeader locks the EXPR header bit layout (§design 1): the four
// legacy kinds keep their old `kind<<14 | count` bit pattern (bit 13 zero), the
// three new kinds set bit 13, pack/unpack round-trips, and a child count above
// 8191 is a hard error.
func TestPackExprHeader(t *testing.T) {
	cases := []struct {
		kind  ExprKind
		count int
		want  uint16
	}{
		{ExprCompare, 0, 0x0000},
		{ExprLogic, 3, 0x4003},
		{ExprNot, 1, 0x8001},
		{ExprEither, 2, 0xC002},
		{ExprExists, 1, 0x2001},
		{ExprCall, 0, 0x6000},
		{ExprCompareX, 0, 0xA000},
		{ExprLogic, exprCountMax, 0x5FFF},
	}
	for _, c := range cases {
		w, err := packExprHeader(c.kind, c.count)
		if err != nil || uint16(w) != c.want {
			t.Fatalf("packExprHeader(%d, %d) = 0x%04X err=%v, want 0x%04X", c.kind, c.count, uint16(w), err, c.want)
		}
		gotKind, gotCount := unpackExprHeader(w)
		if gotKind != c.kind || gotCount != c.count {
			t.Fatalf("unpackExprHeader(0x%04X) = (%d, %d), want (%d, %d)", uint16(w), gotKind, gotCount, c.kind, c.count)
		}
	}
	if _, err := packExprHeader(ExprLogic, exprCountMax+1); err == nil {
		t.Error("packExprHeader(count 8192): expected error")
	}
}
