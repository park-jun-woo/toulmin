//ff:func feature=engine type=engine control=sequence
//ff:what TestLoadGraph_MissingFunc — tests LoadGraph error for missing function
package toulmin

import (
	"testing"
)

func TestLoadGraph_MissingFunc(t *testing.T) {
	def := GraphDef{
		Graph: "missing",
		Rules: []GraphRuleDef{
			{Name: "unknown", Role: "rule"},
		},
	}

	_, err := LoadGraph(def, map[string]any{}, nil)
	if err == nil {
		t.Error("expected error for missing function")
	}
}
