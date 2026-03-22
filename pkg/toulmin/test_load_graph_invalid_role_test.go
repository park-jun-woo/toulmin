//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_InvalidRole — tests LoadGraph error for invalid role
package toulmin

import (
	"testing"
)

func TestLoadGraph_InvalidRole(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }

	def := GraphDef{
		Graph: "bad-role",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "invalid"},
		},
	}

	_, err := LoadGraph(def, map[string]any{"W": fn}, nil)
	if err == nil {
		t.Error("expected error for invalid role")
	}
}
