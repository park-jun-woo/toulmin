//ff:func feature=codegen type=codegen control=sequence
//ff:what TestGenerateGraphGofmtValid — tests that generated code passes gofmt
package codegen

import (
	"go/format"
	"testing"

	"github.com/park-jun-woo/toulmin/internal/graphdef"
)

func TestGenerateGraphGofmtValid(t *testing.T) {
	def := &graphdef.GraphDef{
		Graph: "test",
		Rules: []graphdef.RuleDef{
			{Name: "A", Role: "warrant", Qualifier: float64Ptr(1.0)},
			{Name: "B", Role: "defeater", Qualifier: float64Ptr(0.5)},
		},
		Defeats: []graphdef.EdgeDef{{From: "B", To: "A"}},
	}
	code, err := GenerateGraph("pkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := format.Source([]byte(code)); err != nil {
		t.Errorf("generated code is not gofmt-valid: %v", err)
	}
}
