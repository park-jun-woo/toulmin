//ff:func feature=graph type=validator control=sequence
//ff:what TestValidateGraphDefEmptyName — tests empty graph name fails
package toulmin

import "testing"

func TestValidateGraphDefEmptyName(t *testing.T) {
	def := GraphDef{Graph: ""}
	if err := ValidateGraphDef(def); err == nil {
		t.Fatal("expected error for empty graph name")
	}
}
