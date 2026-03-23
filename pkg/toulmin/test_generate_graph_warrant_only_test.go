//ff:func feature=codegen type=codegen control=sequence
//ff:what TestGenerateGraphWarrantOnly — tests code generation with warrant-only graph
package toulmin

import (
	"strings"
	"testing"
)

func TestGenerateGraphWarrantOnly(t *testing.T) {
	def := &GraphDef{
		Graph: "example",
		Rules: []GraphRuleDef{
			{Name: "IsAdult", Role: "warrant", Qualifier: 1.0},
		},
	}
	code, err := GenerateGraph("mypkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(code, "package mypkg") {
		t.Error("missing package declaration")
	}
	if !strings.Contains(code, "g.Warrant(IsAdult") {
		t.Error("missing Warrant call")
	}
	if !strings.Contains(code, "func() *toulmin.Graph") {
		t.Error("missing IIFE wrapper")
	}
}
