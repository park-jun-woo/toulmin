//ff:func feature=engine type=engine control=sequence
//ff:what TestWarrantWithDefeater — tests verdict when defeater attacks warrant
package toulmin

import (
	"testing"
)

func TestWarrantWithDefeater(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(ctx Context, backing Backing) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(ctx Context, backing Backing) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}
