//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphDefCycle — tests circular defeat is detected
package toulmin

import "testing"

func TestValidateGraphDefCycle(t *testing.T) {
	def := GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "A", Role: "warrant"},
			{Name: "B", Role: "rebuttal"},
		},
		Defeats: []GraphEdgeDef{
			{From: "A", To: "B"},
			{From: "B", To: "A"},
		},
	}
	if err := ValidateGraphDef(def); err == nil {
		t.Fatal("expected error for cycle")
	}
}
