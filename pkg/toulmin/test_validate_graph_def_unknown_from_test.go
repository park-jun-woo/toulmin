//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphDefUnknownFrom — tests defeat from unknown rule fails
package toulmin

import "testing"

func TestValidateGraphDefUnknownFrom(t *testing.T) {
	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "rule"},
		},
		Defeats: []GraphEdgeDef{
			{From: "unknown", To: "W"},
		},
	}
	if err := ValidateGraphDef(def); err == nil {
		t.Fatal("expected error for unknown from rule")
	}
}
