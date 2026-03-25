//ff:func feature=codegen type=codegen control=sequence
//ff:what TestGenerateGraphGofmtValid — tests that generated code passes gofmt
package toulmin

import (
	"go/format"
	"testing"
)

func TestGenerateGraphGofmtValid(t *testing.T) {
	def := &GraphDef{
		Graph: "test",
		Rules: []GraphRuleDef{
			{Name: "A", Role: "rule", Qualifier: 1.0},
			{Name: "B", Role: "except", Qualifier: 0.5},
		},
		Defeats: []GraphEdgeDef{{From: "B", To: "A"}},
	}
	code, err := GenerateGraph("pkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := format.Source([]byte(code)); err != nil {
		t.Errorf("generated code is not gofmt-valid: %v", err)
	}
}
