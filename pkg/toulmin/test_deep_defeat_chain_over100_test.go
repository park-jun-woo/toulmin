//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestDeepDefeatChainOver100 — tests 150-deep defeat chain produces finite verdict
package toulmin

import (
	"math"
	"testing"
)

func TestDeepDefeatChainOver100(t *testing.T) {
	fn := func(ctx Context, backing Backing) (bool, any) { return true, nil }
	eng := NewEngine()
	eng.Register(RuleMeta{Name: "W", Qualifier: 1.0, Strength: Defeasible, Fn: fn})
	prev := "W"
	for i := 1; i <= 150; i++ {
		name := string(rune('A'+i%26)) + string(rune('0'+i/26))
		eng.Register(RuleMeta{Name: name, Qualifier: 1.0, Strength: Defeater, Defeats: []string{prev}, Fn: fn})
		prev = name
	}
	results, err := eng.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := results[0].Verdict
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("verdict should be finite for deep non-cyclic chain, got %f", v)
	}
}
