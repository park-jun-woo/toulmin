//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateExceptGate — Except gates applicability multiplicatively while Counter stays in the denominator
package toulmin

import (
	"math"
	"testing"
)

func TestGraphEvaluateExceptGate(t *testing.T) {
	build := func(v float64, asExcept bool, n int) *Graph {
		g := NewGraph("gate")
		w := g.Rule(WarrantA)
		for i := 0; i < n; i++ {
			fn := func(ctx Context, specs Specs) (float64, any) { return v, nil }
			var a *Rule
			if asExcept {
				a = g.Except(fn).With(&testSpec{Value: string(rune('a' + i))})
			} else {
				a = g.Counter(fn).With(&testSpec{Value: string(rune('a' + i))})
			}
			a.Attacks(w)
		}
		return g
	}
	cases := []struct {
		name     string
		v        float64
		except   bool
		verdict  float64
		applic   float64
		describe string
	}{
		{"ThreeCounters", 1.0, false, -0.5, 1.0, "S = 3 → raw = 1/4 → -0.5 (unchanged)"},
		{"ThreeFullExcepts", 1.0, true, -1.0, 0.0, "gate = 0³ → a_eff 0 → -1"},
		{"ThreePartialExcepts", 0.5, true, 2*math.Pow(0.25, 3) - 1, math.Pow(0.25, 3), "gate = (1-0.75)³"},
		{"InactiveExcept", -1.0, true, 1.0, 1.0, "inactive exception leaves the gate at 1"},
	}
	for _, c := range cases {
		res, err := build(c.v, c.except, 3).Evaluate(NewContext())
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if len(res) != 1 {
			t.Fatalf("%s: expected 1 result, got %d", c.name, len(res))
		}
		if math.Abs(res[0].Verdict-c.verdict) > 1e-12 || math.Abs(res[0].Applicability-c.applic) > 1e-12 {
			t.Errorf("%s (%s): got verdict %v applicability %v", c.name, c.describe, res[0].Verdict, res[0].Applicability)
		}
	}
}
