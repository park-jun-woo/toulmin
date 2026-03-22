//ff:func feature=codegen type=codegen control=sequence
//ff:what TestGenerateGraphWithDefeat — tests code generation with defeat edges
package codegen

import (
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/internal/graphdef"
)

func TestGenerateGraphWithDefeat(t *testing.T) {
	def := &graphdef.GraphDef{
		Graph: "check",
		Rules: []graphdef.RuleDef{
			{Name: "W", Role: "warrant", Qualifier: float64Ptr(1.0)},
			{Name: "R", Role: "rebuttal", Qualifier: float64Ptr(0.8)},
		},
		Defeats: []graphdef.EdgeDef{{From: "R", To: "W"}},
	}
	code, err := GenerateGraph("pkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(code, "g.Rebuttal(R") {
		t.Error("missing Rebuttal call")
	}
	if !strings.Contains(code, "g.Defeat(r, w)") {
		t.Error("missing Defeat call with variable references")
	}
}
