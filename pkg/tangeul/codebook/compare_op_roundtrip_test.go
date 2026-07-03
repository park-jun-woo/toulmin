//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestCompareOpRoundTrip — all 8 Compare.Op strings map and roundtrip, unknown errors
package codebook

import "testing"

func TestCompareOpRoundTrip(t *testing.T) {
	cases := []struct {
		op   string
		code CompareOp
	}{
		{"is empty", CompareIsEmpty},
		{"is not empty", CompareIsNotEmpty},
		{"equals", CompareEquals},
		{"is in", CompareIsIn},
		{"is greater than", CompareIsGreaterThan},
		{"is less than", CompareIsLessThan},
		{"is at most", CompareIsAtMost},
		{"is at least", CompareIsAtLeast},
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
	if _, err := CompareOpToAST(CompareOp(0xFFFF)); err == nil {
		t.Error("CompareOpToAST(unknown): expected error")
	}
}
