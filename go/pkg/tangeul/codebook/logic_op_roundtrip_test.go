//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestLogicOpRoundTrip — "and"/"or"↔LogicOp map and roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/geul/go/pkg/toulmin"
)

func TestLogicOpRoundTrip(t *testing.T) {
	cases := []struct {
		op   string
		code toulmin.LogicOp
	}{
		{"and", toulmin.LogicAnd},
		{"or", toulmin.LogicOr},
	}
	for _, c := range cases {
		code, err := LogicOpFromAST(c.op)
		if err != nil || code != c.code {
			t.Fatalf("LogicOpFromAST(%q): got 0x%04X err=%v, want 0x%04X", c.op, uint16(code), err, uint16(c.code))
		}
		op, err := LogicOpToAST(c.code)
		if err != nil || op != c.op {
			t.Fatalf("LogicOpToAST(0x%04X): got %q err=%v, want %q", uint16(c.code), op, err, c.op)
		}
	}
	if _, err := LogicOpFromAST("xor"); err == nil {
		t.Error("LogicOpFromAST(unknown): expected error")
	}
	if _, err := LogicOpToAST(toulmin.LogicOp(0xFFFF)); err == nil {
		t.Error("LogicOpToAST(unknown): expected error")
	}
}
