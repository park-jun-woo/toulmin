//ff:func feature=engine type=engine control=sequence
//ff:what TestLazySkipsRebuttalWhenWarrantFalse — tests lazy eval skips rebuttal when warrant is false
package toulmin

import (
	"testing"
)

func TestLazySkipsRebuttalWhenWarrantFalse(t *testing.T) {
	rebuttalCalled := false
	falseWarrant := func(ctx Context, backing Backing) (bool, any) { return false, nil }
	trackedRebuttal := func(ctx Context, backing Backing) (bool, any) {
		rebuttalCalled = true
		return true, nil
	}
	g := NewGraph("test")
	w := g.Rule(falseWarrant)
	r := g.Counter(trackedRebuttal)
	r.Attacks(w)
	results, err := g.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results (warrant false), got %d", len(results))
	}
	if rebuttalCalled {
		t.Error("rebuttal func should not be called when warrant is false")
	}
}
