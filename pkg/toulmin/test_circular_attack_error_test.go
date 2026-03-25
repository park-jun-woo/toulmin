//ff:func feature=engine type=engine control=sequence
//ff:what TestCircularAttackError — tests that circular defeat graph returns error
package toulmin

import (
	"testing"
)

func TestCircularAttackError(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "A", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"B"},
		Fn:      func(ctx Context, backing Backing) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "B", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"A"},
		Fn:      func(ctx Context, backing Backing) (bool, any) { return true, nil },
	})
	_, err := eng.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
}
