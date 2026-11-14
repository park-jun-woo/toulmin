//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCompareOpRoundTrip — all 8 Compare.Op strings map and roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestCompareOpRoundTrip(t *testing.T) {
	cases := []struct {
		op   string
		code toulmin.CompareOp
	}{
		{"is empty", toulmin.CompareIsEmpty},
		{"is not empty", toulmin.CompareIsNotEmpty},
		{"equals", toulmin.CompareEquals},
		{"is in", toulmin.CompareIsIn},
		{"is greater than", toulmin.CompareIsGreaterThan},
		{"is less than", toulmin.CompareIsLessThan},
		{"is at most", toulmin.CompareIsAtMost},
		{"is at least", toulmin.CompareIsAtLeast},
	}
	for _, c := range cases {
		code, err := CompareOpFromAST(c.op)
		if err != nil || code != c.code {
			t.Fatalf("CompareOpFromAST(%q): got 0x%04X err=%v, want 0x%04X", c.op, uint16(code), err, uint16(c.code))
		}
		op, err := CompareOpToAST(c.code)
		if err != nil || op != c.op {
			t.Fatalf("CompareOpToAST(0x%04X): got %q err=%v, want %q", uint16(c.code), op, err, c.op)
		}
	}
	if _, err := CompareOpFromAST("is roughly"); err == nil {
		t.Error("CompareOpFromAST(unknown): expected error")
	}
	if _, err := CompareOpToAST(toulmin.CompareOp(0xFFFF)); err == nil {
		t.Error("CompareOpToAST(unknown): expected error")
	}
}
