//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateReuse — Evaluate on one graph many times, with facts that change from call to call, gives each time exactly what a graph built anew gives for the same facts: nothing of an earlier evaluation is left in the reused state
package toulmin

import (
	"reflect"
	"testing"
)

func TestGraphEvaluateReuse(t *testing.T) {
	read := func(ctx Context, specs Specs) (float64, any) {
		key := specs[0].(*testSpec).Value
		v, _ := ctx.Get(key)
		f, _ := v.(float64)
		return f, key
	}
	of := func(key string) Spec { return &testSpec{Value: key} }
	build := func() *Graph {
		g := NewGraph("reuse")
		w := g.Rule(read).With(of("w")).Qualifier(0.9)
		g.Counter(read).With(of("c")).Attacks(w, 0.7)
		x := g.Except(read).With(of("x"))
		x.Attacks(w)
		g.Counter(read).With(of("d")).Attacks(x)
		g.Rule(read).With(of("v"))
		return g
	}
	values := []float64{-1, -0.5, 0, 0.5, 1}
	kept := build()
	for n := 0; n < 625; n++ {
		ctx, same := NewContext(), NewContext()
		for i, key := range []string{"w", "c", "x", "d"} {
			f := values[n/[]int{1, 5, 25, 125}[i]%5]
			ctx.Set(key, f)
			same.Set(key, f)
		}
		ctx.Set("v", values[n%5])
		same.Set("v", values[n%5])
		got, err := kept.Evaluate(ctx)
		want, werr := build().Evaluate(same)
		if err != nil || werr != nil || !reflect.DeepEqual(got, want) {
			t.Fatalf("facts %d: got %+v (%v), want %+v (%v)", n, got, err, want, werr)
		}
	}
}
