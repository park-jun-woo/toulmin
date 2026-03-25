//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphDefUnknownTo — tests defeat to unknown rule fails
package toulmin

import "testing"

func TestValidateGraphDefUnknownTo(t *testing.T) {
	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "R", Role: "counter"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "unknown"},
		},
	}
	if err := ValidateGraphDef(def); err == nil {
		t.Fatal("expected error for unknown to rule")
	}
}
