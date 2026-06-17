//ff:func feature=engine type=engine control=sequence
//ff:what TestRunRecursiveMethod — Run returns the option-resolution error for unsupported method
package toulmin

import "testing"

func TestRunRecursiveMethod(t *testing.T) {
	g := NewGraph("recursive")
	g.Rule(WarrantA)

	results, view, err := g.Run(NewContext(), EvalOption{Method: Recursive})
	if err == nil {
		t.Fatal("expected error for Recursive method")
	}
	if results != nil || view != nil {
		t.Errorf("on option error Run must return nil results and nil view, got results=%v view=%v", results, view)
	}
}
