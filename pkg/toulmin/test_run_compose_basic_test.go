//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeBasic — an Active node Runs its sub-graph (verified via ctx side effect)
package toulmin

import "testing"

func TestRunComposeBasic(t *testing.T) {
	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	subRule := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	sub := NewGraph("sub")
	sub.Rule(subRule).RunOn(func(self TraceEntry, t Trace) error {
		t.Ctx().Set("sub-ran", true)
		return nil
	})

	parent := NewGraph("parent")
	parent.Rule(active).Run(sub)

	ctx := NewContext()
	if _, _, err := parent.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if v, ok := ctx.Get("sub-ran"); !ok || v != true {
		t.Fatal("Active node did not Run its sub-graph")
	}
}
