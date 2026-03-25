//ff:func feature=codegen type=codegen control=sequence
//ff:what TestGenerateGraphWithDefeat — tests code generation with defeat edges
package toulmin

import (
	"strings"
	"testing"
)

func TestGenerateGraphWithDefeat(t *testing.T) {
	def := &GraphDef{
		Graph: "check",
		Rules: []GraphRuleDef{
			{Name: "W", Role: "rule", Qualifier: 1.0},
			{Name: "R", Role: "counter", Qualifier: 0.8},
		},
		Defeats: []GraphEdgeDef{{From: "R", To: "W"}},
	}
	code, err := GenerateGraph("pkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(code, "g.Counter(R") {
		t.Error("missing Rebuttal call")
	}
	if !strings.Contains(code, "r.Attacks(w)") {
		t.Error("missing Defeat call with variable references")
	}
}
