//ff:func feature=engine type=engine control=sequence
//ff:what TestNilFuncGuard — tests that nil func attacker is ignored
package toulmin

import (
	"testing"
)

func TestNilFuncGuard(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(ctx Context, backing Backing) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "Ghost", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"W"},
		Fn:      nil,
	})
	results, err := eng.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (nil attacker ignored), got %f", results[0].Verdict)
	}
}
