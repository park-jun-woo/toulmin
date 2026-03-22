//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphCycleError — tests that circular defeat graph returns error
package toulmin

import (
	"testing"
)

func TestGraphCycleError(t *testing.T) {
	cycleA := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	cycleB := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	a := g.Warrant(cycleA, nil, 1.0)
	b := g.Rebuttal(cycleB, nil, 1.0)
	g.Defeat(b, a)
	g.Defeat(a, b)
	_, err := g.Evaluate(nil, nil)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
}
