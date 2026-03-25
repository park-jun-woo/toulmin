//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_WithDefeater — tests LoadGraph with defeater neutralizing rebuttal
package toulmin

import (
	"testing"
)

func TestLoadGraph_WithDefeater(t *testing.T) {
	wFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	rFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	dFn := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "defeater-test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "rule"},
			{Name: "R", Role: "counter"},
			{Name: "D", Role: "except"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "W"},
			{From: "D", To: "R"},
		},
	}

	funcs := map[string]any{"W": wFn, "R": rFn, "D": dFn}

	g, err := LoadGraph(def, funcs, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil)
	if len(results) == 0 {
		t.Fatal("expected results")
	}
	// W is attacked by R, but R is defeated by D → W should prevail
	if results[0].Verdict <= 0 {
		t.Errorf("expected positive verdict (defeater neutralizes rebuttal), got %f", results[0].Verdict)
	}
}
