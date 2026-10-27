//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestQuantOpRoundTrip — "any"↔QuantOp map and roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestQuantOpRoundTrip(t *testing.T) {
	cases := []struct {
		op   string
		code toulmin.QuantOp
	}{
		{"any", toulmin.QuantAny},
	}
	for _, c := range cases {
		code, err := QuantOpFromAST(c.op)
		if err != nil || code != c.code {
			t.Fatalf("QuantOpFromAST(%q): got 0x%04X err=%v, want 0x%04X", c.op, uint16(code), err, uint16(c.code))
		}
		op, err := QuantOpToAST(c.code)
		if err != nil || op != c.op {
			t.Fatalf("QuantOpToAST(0x%04X): got %q err=%v, want %q", uint16(c.code), op, err, c.op)
		}
	}
	if _, err := QuantOpFromAST("all"); err == nil {
		t.Error("QuantOpFromAST(unknown): expected error")
	}
	if _, err := QuantOpToAST(toulmin.QuantOp(0xFFFF)); err == nil {
		t.Error("QuantOpToAST(unknown): expected error")
	}
}
