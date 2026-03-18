package graphdef

import (
	"testing"
)

func TestValidateValid(t *testing.T) {
	def := &GraphDef{
		Graph: "test",
		Rules: []RuleDef{{Name: "W"}, {Name: "R"}},
		Defeats: []EdgeDef{{From: "R", To: "W"}},
	}
	if err := Validate(def); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateEmptyGraphName(t *testing.T) {
	def := &GraphDef{
		Graph: "",
		Rules: []RuleDef{{Name: "W"}},
	}
	if err := Validate(def); err == nil {
		t.Fatal("expected error for empty graph name")
	}
}

func TestValidateUnknownFrom(t *testing.T) {
	def := &GraphDef{
		Graph:   "test",
		Rules:   []RuleDef{{Name: "W"}},
		Defeats: []EdgeDef{{From: "Ghost", To: "W"}},
	}
	if err := Validate(def); err == nil {
		t.Fatal("expected error for unknown from rule")
	}
}

func TestValidateUnknownTo(t *testing.T) {
	def := &GraphDef{
		Graph:   "test",
		Rules:   []RuleDef{{Name: "R"}},
		Defeats: []EdgeDef{{From: "R", To: "Ghost"}},
	}
	if err := Validate(def); err == nil {
		t.Fatal("expected error for unknown to rule")
	}
}

func TestValidateCycleDetected(t *testing.T) {
	def := &GraphDef{
		Graph:   "test",
		Rules:   []RuleDef{{Name: "A"}, {Name: "B"}},
		Defeats: []EdgeDef{{From: "A", To: "B"}, {From: "B", To: "A"}},
	}
	if err := Validate(def); err == nil {
		t.Fatal("expected error for cyclic defeat graph")
	}
}

func TestValidateNoCycle(t *testing.T) {
	def := &GraphDef{
		Graph:   "test",
		Rules:   []RuleDef{{Name: "W"}, {Name: "R"}, {Name: "D"}},
		Defeats: []EdgeDef{{From: "R", To: "W"}, {From: "D", To: "R"}},
	}
	if err := Validate(def); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
