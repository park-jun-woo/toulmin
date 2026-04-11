//ff:func feature=engine type=engine control=sequence
//ff:what TestSpecInTrace — tests that spec value appears in trace entry
package toulmin

import (
	"testing"
)

func TestSpecInTrace(t *testing.T) {
	isInRole := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(isInRole).With(&testSpec{Value: "admin"})
	results, err := g.Evaluate(NewContext(), EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := results[0].Trace[0].Specs[0].(*testSpec)
	if !ok || b.Value != "admin" {
		t.Errorf("expected spec with value 'admin', got %v", results[0].Trace[0].Specs[0])
	}
}
