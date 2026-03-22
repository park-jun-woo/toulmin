//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateEmptyGraphName — tests validation fails for empty graph name
package graphdef

import (
	"testing"
)

func TestValidateEmptyGraphName(t *testing.T) {
	def := &GraphDef{
		Graph: "",
		Rules: []RuleDef{{Name: "W"}},
	}
	if err := Validate(def); err == nil {
		t.Fatal("expected error for empty graph name")
	}
}
