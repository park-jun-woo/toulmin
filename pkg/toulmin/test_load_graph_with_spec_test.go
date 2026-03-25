//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_WithSpec — tests LoadGraph with spec parameter
package toulmin

import (
	"testing"
)

func TestLoadGraph_WithSpec(t *testing.T) {
	fn := func(ctx Context, specs Specs) (bool, any) {
		tb, ok := specs[0].(*testSpec)
		return ok && tb.Value == "admin", specs
	}

	def := GraphDef{
		Graph: "spec-test",
		Rules: []GraphRuleDef{
			{Name: "checkRole", Role: "rule", Qualifier: 1.0},
		},
	}

	funcs := map[string]any{"checkRole": fn}
	specMap := map[string][]Spec{"checkRole": {&testSpec{Value: "admin"}}}

	g, err := LoadGraph(def, funcs, specMap)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil)
	if len(results) == 0 || results[0].Verdict <= 0 {
		t.Error("expected positive verdict with admin spec")
	}
}
