//ff:func feature=engine type=engine control=sequence
//ff:what TestRunComposeCtxFlows — ctx flows down so a sub-graph rule reads a value set by the parent
package toulmin

import "testing"

func TestRunComposeCtxFlows(t *testing.T) {
	subSaw := false
	subRule := func(ctx Context, specs Specs) (bool, any) {
		if v, ok := ctx.Get("token"); ok && v == "abc" {
			subSaw = true
		}
		return true, nil
	}
	sub := NewGraph("sub")
	sub.Rule(subRule)

	active := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	parent := NewGraph("parent")
	parent.Rule(active).Run(sub)

	ctx := NewContext()
	ctx.Set("token", "abc")
	if _, _, err := parent.Run(ctx); err != nil {
		t.Fatalf("run error: %v", err)
	}
	if !subSaw {
		t.Error("sub-graph rule did not read the ctx value set by the parent")
	}
}
