//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeDiamond — two Active nodes sharing one sub-graph is a legal DAG (no cycle)
package toulmin

import "testing"

func TestRunComposeDiamond(t *testing.T) {
	subRuns := 0
	subRule := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	shared := NewGraph("shared")
	shared.Rule(subRule).RunOn(func(t Trace) error {
		subRuns++
		return nil
	})

	ruleX := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	ruleY := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent := NewGraph("parent")
	parent.Rule(ruleX).Run(shared)
	parent.Rule(ruleY).Run(shared)

	if _, _, err := parent.Run(NewContext()); err != nil {
		t.Fatalf("diamond DAG must be legal, got error: %v", err)
	}
	if subRuns != 2 {
		t.Errorf("shared sub-graph should Run once per Active node (2), got %d", subRuns)
	}
}
