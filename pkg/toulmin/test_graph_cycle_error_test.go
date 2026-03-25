//ff:func feature=engine type=engine control=sequence
//ff:what TestGraphCycleError — tests that circular defeat graph returns error
package toulmin

import (
	"testing"
)

func TestGraphCycleError(t *testing.T) {
	cycleA := func(ctx Context, backing Backing) (bool, any) { return true, nil }
	cycleB := func(ctx Context, backing Backing) (bool, any) { return true, nil }
	g := NewGraph("test")
	a := g.Rule(cycleA)
	b := g.Counter(cycleB)
	b.Attacks(a)
	a.Attacks(b)
	_, err := g.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
}
