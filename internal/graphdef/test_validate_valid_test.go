//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateValid — tests validation passes for a valid GraphDef
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
