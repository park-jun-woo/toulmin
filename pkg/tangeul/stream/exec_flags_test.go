//ff:func feature=tangl type=util control=iteration dimension=1
//ff:what TestExecFlags — once/certainty combinations set F_ONCE/F_GATE, and an unknown Certainty op errors
package stream

import (
	"testing"

	"github.com/park-jun-woo/toulmin/pkg/tangeul/codebook"
	"github.com/park-jun-woo/toulmin/pkg/tangl/ast"
)

// TestExecFlags walks every combination of Once and Certainty (absent, valid,
// invalid op) that execFlags branches on.
func TestExecFlags(t *testing.T) {
	cases := []struct {
		name     string
		x        ast.Exec
		wantErr  bool
		wantOnce bool
		wantGate bool
		wantPct  uint8
	}{
		{name: "neither", x: ast.Exec{}, wantOnce: false, wantGate: false},
		{name: "onceOnly", x: ast.Exec{Once: true}, wantOnce: true, wantGate: false},
		{
			name:     "certaintyValid",
			x:        ast.Exec{Certainty: &ast.Certainty{Op: "at least", Percent: 75}},
			wantOnce: false, wantGate: true, wantPct: 75,
		},
		{
			name:     "onceAndCertainty",
			x:        ast.Exec{Once: true, Certainty: &ast.Certainty{Op: "above", Percent: 50}},
			wantOnce: true, wantGate: true, wantPct: 50,
		},
		{
			name:    "certaintyInvalidOp",
			x:       ast.Exec{Certainty: &ast.Certainty{Op: "bogus", Percent: 10}},
			wantErr: true,
		},
	}
	for _, c := range cases {
		ep, err := execFlags(c.x)
		if c.wantErr && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
		if c.wantErr {
			continue
		}
		if err != nil {
			t.Fatalf("%s: execFlags: %v", c.name, err)
		}
		if got := ep.Flags&codebook.EdgeFlagOnce != 0; got != c.wantOnce {
			t.Errorf("%s: F_ONCE = %v, want %v", c.name, got, c.wantOnce)
		}
		gateSet := ep.Flags&codebook.EdgeFlagGate != 0
		if gateSet != c.wantGate {
			t.Errorf("%s: F_GATE = %v, want %v", c.name, gateSet, c.wantGate)
		}
		if !c.wantGate && ep.Gate != nil {
			t.Errorf("%s: Gate = %v, want nil", c.name, ep.Gate)
		}
		if c.wantGate && ep.Gate == nil {
			t.Fatalf("%s: Gate = nil, want set", c.name)
		}
		if c.wantGate && ep.Gate.Percent != c.wantPct {
			t.Errorf("%s: Gate.Percent = %d, want %d", c.name, ep.Gate.Percent, c.wantPct)
		}
	}
}
