//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAssignPrefix — assignPrefix records TEXTs, recovers the Rules condition first, and short-circuits on error
package stream

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"

	"github.com/geul-org/geul/go/pkg/core"
	"github.com/park-jun-woo/toulmin/go/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/go/pkg/tangl/ast"
)

// TestAssignPrefix runs an empty prefix, an all-TEXT prefix, a Rules prefix whose
// first record is the condition, and a malformed record.
func TestAssignPrefix(t *testing.T) {
	text := func(id toulmin.TID, s string) rawPacket {
		ws, _ := toulmin.EncodeText(toulmin.TextPacket{TID: id, Bytes: []byte(s)}, testWire)
		return rawOf(t, ws)
	}
	ew, _ := packet.EncodeExpr(packet.ExprPacket{ID: 9, Cond: ast.Compare{Field: "n", Op: "equals", Value: "0"}}, testWire)
	bad := rawPacket{Class: classPayload, Rec: core.Record{Kind: core.KindEdge, Codec: core.CodecIssuer1Word, Issuer: []core.Word{0}, Format: toulmin.FormatText}}
	cases := []struct {
		prefix   []rawPacket
		isRules  bool
		wantErr  bool
		wantCond bool
		texts    int
	}{
		{nil, false, false, false, 0},
		{[]rawPacket{text(1, "a"), text(2, "b")}, false, false, false, 2},
		{[]rawPacket{rawOf(t, ew), text(3, "c")}, true, false, true, 1},
		{[]rawPacket{bad}, false, true, false, 0},
	}
	for i, c := range cases {
		d := newDec(testWire)
		got, err := assignPrefix(d, c.prefix)
		if (err != nil) != c.wantErr || (got != nil) != c.wantCond || len(d.texts) != c.texts {
			t.Fatalf("case %d: cond %v err %v texts %d", i, got, err, len(d.texts))
		}
	}
}
