//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_MissingDefeatRef — tests LoadGraph error for missing defeat reference
package toulmin

import (
	"testing"
)

func TestLoadGraph_MissingDefeatRef(t *testing.T) {
	fn := func(c any, g any, b Backing) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "bad-edge",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
		},
		Defeats: []GraphEdgeDef{
			{From: "ghost", To: "W"},
		},
	}

	_, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err == nil {
		t.Error("expected error for missing defeat reference")
	}
}
