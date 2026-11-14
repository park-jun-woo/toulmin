//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleRun — (*Rule).Run wires a non-nil sub-graph and panics on nil
package toulmin

import "testing"

// TestRuleRun covers both branches of (*Rule).Run: the normal path sets RunGraph and
// returns the receiver for chaining, and a nil sub-graph is a registration error that panics.
func TestRuleRun(t *testing.T) {
	// Distinct closures: two rules in one graph must be different function values.
	f1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	f2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	// normal: Run sets the node's RunGraph edge and returns the same *Rule.
	sub := NewGraph("sub")
	sub.Rule(f1)
	parent := NewGraph("parent")
	r := parent.Rule(f1)
	if got := r.Run(sub); got != r {
		t.Errorf("Run must return the receiver for chaining, got %v want %v", got, r)
	}
	if parent.rules[r.idx].RunGraph != sub {
		t.Errorf("Run must set RunGraph to the sub-graph, got %v", parent.rules[r.idx].RunGraph)
	}

	// nil sub-graph: registration error -> panic.
	func() {
		defer func() {
			if rec := recover(); rec == nil {
				t.Error("Run(nil) must panic")
			}
		}()
		parent.Rule(f2).Run(nil)
	}()
}
