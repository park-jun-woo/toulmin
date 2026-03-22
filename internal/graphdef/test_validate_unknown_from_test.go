//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateUnknownFrom — tests validation fails for unknown 'from' rule in defeat edge
package graphdef

import (
	"testing"
)

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
