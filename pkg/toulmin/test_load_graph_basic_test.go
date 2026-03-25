//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_Basic — tests basic LoadGraph with warrant and rebuttal
package toulmin

import (
	"testing"
)

func TestLoadGraph_Basic(t *testing.T) {
	wFn := func(c any, g any, b Backing) (bool, any) { return true, nil }
	rFn := func(c any, g any, b Backing) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
			{Name: "R", Role: "rebuttal"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "W"},
		},
	}

	funcs := map[string]any{"W": wFn, "R": rFn}

	g, err := LoadGraph(def, funcs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("evaluate error: %v", err)
	}
	if len(results) == 0 {
		t.Fatal("expected results")
	}
}
