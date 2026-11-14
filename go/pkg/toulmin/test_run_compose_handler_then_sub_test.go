//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeHandlerThenSub — the node's RunOn handler fires before its sub-graph Runs
package toulmin

import "testing"

func TestRunComposeHandlerThenSub(t *testing.T) {
	var order []string
	subRule := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	sub := NewGraph("sub")
	sub.Rule(subRule).RunOn(func(self TraceEntry, t Trace) error {
		order = append(order, "sub")
		return nil
	})

	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent := NewGraph("parent")
	parent.Rule(active).RunOn(func(self TraceEntry, t Trace) error {
		order = append(order, "handler")
		return nil
	}).Run(sub)

	if _, _, err := parent.Run(NewContext()); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if len(order) != 2 || order[0] != "handler" || order[1] != "sub" {
		t.Errorf("expected [handler sub], got %v", order)
	}
}
