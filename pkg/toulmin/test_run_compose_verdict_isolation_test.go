//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeVerdictIsolation — a sub-graph's verdict does not leak into the parent's results
package toulmin

import "testing"

func TestRunComposeVerdictIsolation(t *testing.T) {
	// Sub-graph whose warrant is Defeated (verdict 0) — must not affect the parent.
	subRan := false
	subWarrant := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	subAttacker := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	sub := NewGraph("sub")
	sw := sub.Rule(subWarrant).OnDefeated(func(ctx Context, ev NodeEvent, view RunView) error {
		subRan = true
		return nil
	})
	sub.Except(subAttacker).Attacks(sw)

	parent := NewGraph("parent")
	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent.Rule(active).Run(sub)

	results, _, err := parent.Run(NewContext())
	if err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !subRan {
		t.Fatal("sub-graph did not Run")
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 parent result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("parent verdict must stay 1.0 regardless of sub-graph verdict, got %f", results[0].Verdict)
	}
}
