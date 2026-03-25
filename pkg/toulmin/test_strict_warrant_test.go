//ff:func feature=engine type=engine control=sequence
//ff:what TestStrictWarrant — tests that strict warrant ignores attacks
package toulmin

import (
	"testing"
)

func TestStrictWarrant(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Strict,
		Fn: func(ctx Context, specs Specs) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(ctx Context, specs Specs) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (strict rejects attack), got %f", results[0].Verdict)
	}
}
