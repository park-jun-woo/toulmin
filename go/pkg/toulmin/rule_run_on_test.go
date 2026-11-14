//ff:func feature=engine type=engine control=sequence
//ff:what TestRuleRunOn — tests Rule.RunOn wires a NodeHandler and returns the receiver for chaining
package toulmin

import "testing"

func TestRuleRunOn(t *testing.T) {
	f := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	h := func(self TraceEntry, tr Trace) error { return nil }

	g := NewGraph("test")
	r := g.Rule(f)
	got := r.RunOn(h)
	if got != r {
		t.Errorf("RunOn must return the receiver for chaining, got %v want %v", got, r)
	}
	if g.rules[r.idx].RunOn == nil {
		t.Errorf("RunOn must set the node's RunOn handler")
	}
}
