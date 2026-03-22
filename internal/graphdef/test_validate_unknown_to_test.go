//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateUnknownTo — tests validation fails for unknown 'to' rule in defeat edge
package graphdef

import (
	"testing"
)

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
