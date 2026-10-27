//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestGateOpRoundTrip — Certainty.Op↔GateOp for all 4 ops, roundtrip, unknown errors
package codebook

import (
	"testing"

	"github.com/geul-org/gogeul/pkg/toulmin"
)

func TestGateOpRoundTrip(t *testing.T) {
	cases := []struct {
		op   string
		gate toulmin.GateOp
	}{
		{"at least", toulmin.GateAtLeast},
		{"above", toulmin.GateAbove},
		{"less than", toulmin.GateLessThan},
		{"at most", toulmin.GateAtMost},
	}
	for _, c := range cases {
		gate, err := GateOpFromAST(c.op)
		if err != nil || gate != c.gate {
			t.Fatalf("GateOpFromAST(%q): got %d err=%v, want %d", c.op, uint8(gate), err, uint8(c.gate))
		}
		op, err := GateOpToAST(c.gate)
		if err != nil || op != c.op {
			t.Fatalf("GateOpToAST(%d): got %q err=%v, want %q", uint8(c.gate), op, err, c.op)
		}
	}
	if _, err := GateOpFromAST("maybe"); err == nil {
		t.Error("GateOpFromAST(unknown): expected error")
	}
	if _, err := GateOpToAST(toulmin.GateOp(99)); err == nil {
		t.Error("GateOpToAST(unknown): expected error")
	}
}
