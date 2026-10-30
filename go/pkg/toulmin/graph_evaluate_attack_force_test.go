//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateAttackForce — attack force scales both the rebuttal denominator and the exception gate
package toulmin

import (
	"math"
	"testing"
)

func TestGraphEvaluateAttackForce(t *testing.T) {
	cases := []struct {
		name    string
		except  bool
		force   float64
		verdict float64
	}{
		{"CounterForceZero", false, 0.0, 1.0},
		{"CounterForceHalf", false, 0.5, 1.0 / 3.0}, // S = 0.5 → raw 2/3
		{"CounterForceOne", false, 1.0, 0.0},
		{"ExceptForceZero", true, 0.0, 1.0},
		{"ExceptForceHalf", true, 0.5, 0.0}, // a_eff = 0.5
		{"ExceptForceQuarter", true, 0.25, 0.5},
	}
	for i := 0; i < 2*len(cases); i++ {
		c := cases[i/2]
		opt := EvalOption{Trace: i%2 == 1}
		g := NewGraph("force")
		w := g.Rule(WarrantA)
		a := g.Counter(RebuttalB)
		if c.except {
			a = g.Except(DefeaterC)
		}
		a.Attacks(w, c.force)
		res, err := g.Evaluate(NewContext(), opt)
		if err != nil {
			t.Fatalf("%s: %v", c.name, err)
		}
		if math.Abs(res[0].Verdict-c.verdict) > 1e-12 {
			t.Errorf("%s (trace=%v): expected %v, got %v", c.name, opt.Trace, c.verdict, res[0].Verdict)
		}
	}
}
