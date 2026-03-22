//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateNoCycle — tests validation passes for acyclic defeat graph
package graphdef

import (
	"testing"
)

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
