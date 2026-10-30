//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateBoolCompat — bool rules give exactly the same verdicts as ±1 float64 rules and the legacy formula
package toulmin

import "testing"

func TestGraphEvaluateBoolCompat(t *testing.T) {
	boolOf := func(b bool) func(Context, Specs) (bool, any) {
		return func(ctx Context, specs Specs) (bool, any) { return b, nil }
	}
	floatOf := func(b bool) func(Context, Specs) (float64, any) {
		v := -1.0
		if b {
			v = 1.0
		}
		return func(ctx Context, specs Specs) (float64, any) { return v, nil }
	}
	// w ← c1 (q 0.7), c2 (q 0.4, inactive), c3 (q 0.3); c1 ← c4 (q 0.9)
	build := func(useBool bool) *Graph {
		g := NewGraph("compat")
		reg := func(add func(any) *Rule, b bool, tag string) *Rule {
			if useBool {
				return add(boolOf(b)).With(&testSpec{Value: tag})
			}
			return add(floatOf(b)).With(&testSpec{Value: tag})
		}
		w := reg(g.Rule, true, "w").Qualifier(0.8)
		c1 := reg(g.Counter, true, "c1").Qualifier(0.7)
		c2 := reg(g.Counter, false, "c2").Qualifier(0.4)
		c3 := reg(g.Counter, true, "c3").Qualifier(0.3)
		c4 := reg(g.Counter, true, "c4").Qualifier(0.9)
		c1.Attacks(w)
		c2.Attacks(w)
		c3.Attacks(w)
		c4.Attacks(c1)
		return g
	}
	// legacy formula: raw = q/(1+Σ(v_j+1)/2), v = 2raw−1
	vc4 := 2*0.9 - 1
	vc1 := 2*(0.7/(1.0+(vc4+1.0)/2.0)) - 1
	vc3 := 2*0.3 - 1
	legacy := 2*(0.8/(1.0+((vc1+1.0)/2.0+0.0+(vc3+1.0)/2.0))) - 1
	for _, opt := range []EvalOption{{}, {Trace: true}} {
		rb, err := build(true).Evaluate(NewContext(), opt)
		if err != nil {
			t.Fatal(err)
		}
		rf, err := build(false).Evaluate(NewContext(), opt)
		if err != nil {
			t.Fatal(err)
		}
		if rb[0].Verdict != legacy || rf[0].Verdict != legacy {
			t.Errorf("trace=%v: bool %v / float %v must equal legacy %v exactly", opt.Trace, rb[0].Verdict, rf[0].Verdict, legacy)
		}
	}
}
