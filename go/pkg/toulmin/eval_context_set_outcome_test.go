//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestEvalContextSetOutcome — setOutcome stores v and evidence, and active is exactly v > -1
package toulmin

import "testing"

func TestEvalContextSetOutcome(t *testing.T) {
	cases := []struct {
		v      float64
		active bool
	}{
		{-1.0, false},
		{-0.999, true},
		{0.0, true},
		{1.0, true},
	}
	for _, c := range cases {
		ec := &evalContext{
			ruleVerdict: map[string]float64{},
			active:      map[string]bool{},
			evidence:    map[string]any{},
		}
		ec.setOutcome("a", c.v, "ev")
		if ec.ruleVerdict["a"] != c.v {
			t.Errorf("v=%v: ruleVerdict = %v", c.v, ec.ruleVerdict["a"])
		}
		if ec.active["a"] != c.active {
			t.Errorf("v=%v: active = %v, want %v", c.v, ec.active["a"], c.active)
		}
		if ec.evidence["a"] != "ev" {
			t.Errorf("v=%v: evidence = %v", c.v, ec.evidence["a"])
		}
	}
}
