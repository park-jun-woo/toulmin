//ff:func feature=engine type=engine control=sequence
//ff:what TestCompensation — tests verdict with compensation chain (W ← D1 ← D2)
package toulmin

import (
	"math"
	"testing"
)

func TestCompensation(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D1", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D2", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"D1"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}
