//ff:func feature=analyzer type=analyzer control=sequence
//ff:what TestExtractDefeatsMethonChain — tests defeat extraction from method chain style Graph builder
package analyzer

import (
	"testing"
)

func TestExtractDefeatsMethonChain(t *testing.T) {
	path := writeGoFile(t, `package example

import "github.com/park-jun-woo/toulmin/pkg/toulmin"

func W(c, g any) (bool, any) { return true, nil }
func R(c, g any) (bool, any) { return true, nil }

var g = toulmin.NewGraph("check").
	Warrant(W, 1.0).
	Rebuttal(R, 1.0).
	Defeat(R, W)
`)
	graphs, err := ExtractDefeats(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(graphs) != 1 {
		t.Fatalf("expected 1 graph, got %d", len(graphs))
	}
	dg := graphs[0]
	if dg.Graph != "check" {
		t.Errorf("expected graph name 'check', got %q", dg.Graph)
	}
	if len(dg.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(dg.Rules))
	}
	if len(dg.Defeats) != 1 || dg.Defeats[0].From != "R" || dg.Defeats[0].To != "W" {
		t.Errorf("expected Defeats=[{R W}], got %v", dg.Defeats)
	}
}
