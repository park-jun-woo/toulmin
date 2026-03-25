//ff:func feature=engine type=engine control=sequence
//ff:what TestWarrantOnly — tests verdict for single warrant node
package toulmin

import (
	"testing"
)

func TestWarrantOnly(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(ctx Context, specs Specs) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
