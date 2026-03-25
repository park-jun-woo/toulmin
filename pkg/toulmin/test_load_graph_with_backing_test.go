//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_WithBacking — tests LoadGraph with backing parameter
package toulmin

import (
	"testing"
)

func TestLoadGraph_WithBacking(t *testing.T) {
	fn := func(ctx Context, backing Backing) (bool, any) {
		tb, ok := backing.(*testBacking)
		return ok && tb.Value == "admin", backing
	}

	def := GraphDef{
		Graph: "backing-test",
		Rules: []GraphRuleDef{
			{Name: "checkRole", Role: "rule", Qualifier: 1.0},
		},
	}

	funcs := map[string]any{"checkRole": fn}
	backings := map[string]Backing{"checkRole": &testBacking{Value: "admin"}}

	g, err := LoadGraph(def, funcs, backings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil)
	if len(results) == 0 || results[0].Verdict <= 0 {
		t.Error("expected positive verdict with admin backing")
	}
}
