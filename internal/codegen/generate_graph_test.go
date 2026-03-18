package codegen

import (
	"go/format"
	"strings"
	"testing"

	"github.com/park-jun-woo/toulmin/internal/graphdef"
)

func float64Ptr(v float64) *float64 { return &v }

func TestGenerateGraphWarrantOnly(t *testing.T) {
	def := &graphdef.GraphDef{
		Graph: "example",
		Rules: []graphdef.RuleDef{
			{Name: "IsAdult", Role: "warrant", Qualifier: float64Ptr(1.0)},
		},
	}
	code, err := GenerateGraph("mypkg", def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(code, "package mypkg") {
		t.Error("missing package declaration")
	}
	if !strings.Contains(code, "Warrant(IsAdult") {
		t.Error("missing Warrant call")
	}
}

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
	if !strings.Contains(code, "Rebuttal(R") {
		t.Error("missing Rebuttal call")
	}
	if !strings.Contains(code, "Defeat(R, W)") {
		t.Error("missing Defeat call")
	}
}

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
