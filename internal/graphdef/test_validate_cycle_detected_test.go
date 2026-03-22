//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateCycleDetected — tests validation fails for cyclic defeat graph
package graphdef

import (
	"testing"
)

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
