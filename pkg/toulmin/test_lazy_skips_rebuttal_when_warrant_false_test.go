//ff:func feature=engine type=engine control=sequence
//ff:what TestLazySkipsRebuttalWhenWarrantFalse — tests lazy eval skips rebuttal when warrant is false
package toulmin

import (
	"testing"
)

func TestLazySkipsRebuttalWhenWarrantFalse(t *testing.T) {
	rebuttalCalled := false
	falseWarrant := func(claim any, ground any, backing Backing) (bool, any) { return false, nil }
	trackedRebuttal := func(claim any, ground any, backing Backing) (bool, any) {
		rebuttalCalled = true
		return true, nil
	}
	g := NewGraph("test")
	w := g.Warrant(falseWarrant, nil, 1.0)
	r := g.Rebuttal(trackedRebuttal, nil, 1.0)
	g.Defeat(r, w)
	results, err := g.Evaluate(nil, nil)
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
