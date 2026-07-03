//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestAssignPrefix — walks a statement's prefix payloads, recovering a Rules EXPR by position, or errors
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/packet"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestAssignPrefix exercises assignPrefix's loop across its zero-iteration case
// (empty prefix, nil cond), an all-Text prefix (isRules false, every payload
// issues a Text TID and cond stays nil), an isRules prefix whose first payload
// is the EXPR condition and whose remaining payloads are Text, and a decode
// error from a malformed payload that short-circuits the loop.
func TestAssignPrefix(t *testing.T) {
	textPacket := func(s string) rawPacket {
		return rawPacket{Class: 2, Words: packet.EncodeText(packet.TextPacket{Bytes: []byte(s)})}
	}
	cond := ast.Compare{Field: "n", Op: "equals", Value: "0"}
	exprWords, err := packet.EncodeExpr(packet.ExprPacket{Cond: cond})
	if err != nil {
		t.Fatalf("EncodeExpr fixture: %v", err)
	}
	exprPacket := rawPacket{Class: 2, Words: exprWords}

	cases := []struct {
		name     string
		prefix   []rawPacket
		isRules  bool
		wantErr  bool
		wantCond bool
		wantTIDs int
	}{
		{"empty prefix, zero iterations", nil, false, false, false, 0},
		{"all-text prefix", []rawPacket{textPacket("a"), textPacket("b")}, false, false, false, 2},
		{"rules prefix: expr then text", []rawPacket{exprPacket, textPacket("c")}, true, false, true, 1},
		{"decode error short-circuits", []rawPacket{{Class: 2, Words: nil}}, false, true, false, 0},
	}

	run := func(t *testing.T, i int) {
		c := cases[i]
		d := &dec{texts: map[uint16]string{}}
		got, err := assignPrefix(d, c.prefix, c.isRules)
		if c.wantErr {
			if err == nil {
				t.Fatalf("want error, got nil")
			}
			return
		}
		if err != nil {
			t.Fatalf("assignPrefix: %v", err)
		}
		if c.wantCond && got == nil {
			t.Errorf("want non-nil cond, got nil")
		}
		if !c.wantCond && got != nil {
			t.Errorf("want nil cond, got %v", got)
		}
		if len(d.texts) != c.wantTIDs {
			t.Errorf("texts len = %d, want %d", len(d.texts), c.wantTIDs)
		}
	}

	for i := range cases {
		t.Run(cases[i].name, func(t *testing.T) { run(t, i) })
	}
}
