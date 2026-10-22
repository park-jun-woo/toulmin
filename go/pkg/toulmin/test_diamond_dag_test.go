//ff:func feature=engine type=engine control=sequence
//ff:what TestDiamondDAG — tests verdict memoization on diamond DAG
package toulmin

import (
	"testing"
)

func TestDiamondDAG(t *testing.T) {
	w := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	r1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	r2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	d := func(ctx Context, specs Specs) (bool, any) { return true, nil }

	g := NewGraph("test")
	warrant := g.Rule(w)
	rebuttal1 := g.Counter(r1)
	rebuttal2 := g.Counter(r2)
	defeater := g.Except(d)

	rebuttal1.Attacks(warrant)
	rebuttal2.Attacks(warrant)
	defeater.Attacks(rebuttal1)
	defeater.Attacks(rebuttal2)

	results, err := g.Evaluate(NewContext())
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}
