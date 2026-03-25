//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_WithBacking — tests LoadGraph with backing parameter
package toulmin

import (
	"testing"
)

func TestLoadGraph_WithBacking(t *testing.T) {
	fn := func(c any, g any, b Backing) (bool, any) {
		tb, ok := b.(*testBacking)
		return ok && tb.Value == "admin", b
	}

	def := GraphDef{
		Graph: "backing-test",
		Rules: []GraphRuleDef{
			{Name: "checkRole", Role: "warrant", Qualifier: 1.0},
		},
	}

	funcs := map[string]any{"checkRole": fn}
	backings := map[string]Backing{"checkRole": &testBacking{Value: "admin"}}

	g, err := LoadGraph(def, funcs, backings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil, nil)
	if len(results) == 0 || results[0].Verdict <= 0 {
		t.Error("expected positive verdict with admin backing")
	}
}
