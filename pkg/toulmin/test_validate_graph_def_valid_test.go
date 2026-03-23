//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphDefValid — tests valid GraphDef passes validation
package toulmin

import "testing"

func TestValidateGraphDefValid(t *testing.T) {
	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "warrant"},
			{Name: "R", Role: "rebuttal"},
		},
		Defeats: []GraphEdgeDef{
			{From: "R", To: "W"},
		},
	}
	if err := ValidateGraphDef(def); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
