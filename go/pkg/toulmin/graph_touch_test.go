//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphTouch — every change to a graph after an evaluation (Rule, Counter, Except, With, Qualifier, Level, Attacks, Supersedes) drops its plan, and the next evaluation gives what a graph built that way from the start gives
package toulmin

import (
	"reflect"
	"testing"
)

func TestGraphTouch(t *testing.T) {
	yes := func(ctx Context, specs Specs) (float64, any) { return 1.0, "yes" }
	half := func(ctx Context, specs Specs) (float64, any) { return 0.0, "half" }
	low := func(ctx Context, specs Specs) (float64, any) { return -0.5, "low" }
	some := func(ctx Context, specs Specs) (float64, any) { return 0.5, len(specs) }
	steps := []func(g *Graph, r map[string]*Rule){
		func(g *Graph, r map[string]*Rule) { r["w"] = g.Rule(yes) },
		func(g *Graph, r map[string]*Rule) { r["c"] = g.Counter(half) },
		func(g *Graph, r map[string]*Rule) { r["c"].Attacks(r["w"]) },
		func(g *Graph, r map[string]*Rule) { r["x"] = g.Except(some) },
		func(g *Graph, r map[string]*Rule) { r["x"].Attacks(r["w"], 0.5) },
		func(g *Graph, r map[string]*Rule) { r["x"].With(&testSpec{Value: "a"}) },
		func(g *Graph, r map[string]*Rule) { r["w"].Qualifier(0.8) },
		func(g *Graph, r map[string]*Rule) { r["c"].Level(Review) },
		func(g *Graph, r map[string]*Rule) { r["s"] = g.Rule(low) },
		func(g *Graph, r map[string]*Rule) { r["s"].Supersedes(r["w"]) },
	}
	kept, rules := NewGraph("touch"), map[string]*Rule{}
	for n := range steps {
		steps[n](kept, rules)
		if kept.compiled.Load() != nil {
			t.Fatalf("step %d: the plan was not dropped", n)
		}
		got, err := kept.Evaluate(NewContext(), EvalOption{Trace: true})
		if err != nil || kept.compiled.Load() == nil {
			t.Fatalf("step %d: %v", n, err)
		}
		fresh, again := NewGraph("touch"), map[string]*Rule{}
		for _, step := range steps[:n+1] {
			step(fresh, again)
		}
		want, err := fresh.Evaluate(NewContext(), EvalOption{Trace: true})
		if err != nil || len(want) == 0 || !reflect.DeepEqual(got, want) {
			t.Fatalf("step %d: got %+v, want %+v (%v)", n, got, want, err)
		}
	}
}
