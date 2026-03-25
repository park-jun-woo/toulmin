//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestDeepDefeatChain — tests 11-deep defeat chain produces finite verdict
package toulmin

import (
	"math"
	"testing"
)

func TestDeepDefeatChain(t *testing.T) {
	fn := func(ctx Context, backing Backing) (bool, any) { return true, nil }
	eng := NewEngine()
	eng.Register(RuleMeta{Name: "W", Qualifier: 1.0, Strength: Defeasible, Fn: fn})
	prev := "W"
	for i := 1; i <= 11; i++ {
		name := "D" + string(rune('0'+i))
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
		t.Errorf("verdict should be finite, got %f", v)
	}
}
