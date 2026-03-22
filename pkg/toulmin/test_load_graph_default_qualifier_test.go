//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_DefaultQualifier — tests LoadGraph uses default qualifier 1.0 when 0
package toulmin

import (
	"testing"
)

func TestLoadGraph_DefaultQualifier(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "default-q",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"}, // Qualifier=0 → default 1.0
		},
	}

	g, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	results, _ := g.Evaluate(nil, nil)
	if len(results) == 0 || results[0].Verdict != 1.0 {
		t.Errorf("expected verdict 1.0 with default qualifier, got %v", results)
	}
}
