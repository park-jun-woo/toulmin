//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluatePartialVerdict — float64 rules: partial verdicts, Applicability, and lazy skip only at v = -1
package toulmin

import (
	"math"
	"testing"
)

func TestGraphEvaluatePartialVerdict(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"HalfApplies", func(t *testing.T) {
			// v = 0 means "half applies" (NOT "not applicable"): a = 0.5 → verdict 0.
			g := NewGraph("partial")
			g.Rule(func(ctx Context, specs Specs) (float64, any) { return 0.0, "half" })
			res, err := g.Evaluate(NewContext())
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 1 || res[0].Verdict != 0.0 || res[0].Applicability != 0.5 || res[0].Evidence != "half" {
				t.Fatalf("expected one result (verdict 0, applicability 0.5, evidence half), got %+v", res)
			}
		}},
		{"QualifierScales", func(t *testing.T) {
			// v = 0.5 → a = 0.75; q = 0.8 → raw = 0.6 → verdict 0.2
			g := NewGraph("partial-q")
			g.Rule(func(ctx Context, specs Specs) (float64, any) { return 0.5, nil }).Qualifier(0.8)
			res, err := g.Evaluate(NewContext())
			if err != nil {
				t.Fatal(err)
			}
			if math.Abs(res[0].Verdict-0.2) > 1e-12 || res[0].Applicability != 0.75 {
				t.Fatalf("expected verdict 0.2 / applicability 0.75, got %+v", res[0])
			}
		}},
		{"MinusOneIsNotApplicableAndLazy", func(t *testing.T) {
			attacked := false
			g := NewGraph("lazy")
			w := g.Rule(func(ctx Context, specs Specs) (float64, any) { return -1.0, nil })
			c := g.Counter(func(ctx Context, specs Specs) (float64, any) { attacked = true; return 1.0, nil })
			c.Attacks(w)
			res, err := g.Evaluate(NewContext())
			if err != nil {
				t.Fatal(err)
			}
			if len(res) != 0 {
				t.Fatalf("v = -1 must be inactive (no result), got %+v", res)
			}
			if attacked {
				t.Fatal("attackers must not be evaluated when the rule returns -1")
			}
		}},
		{"BarelyAppliesEvaluatesAttackers", func(t *testing.T) {
			attacked := false
			g := NewGraph("barely")
			w := g.Rule(func(ctx Context, specs Specs) (float64, any) { return -0.5, nil })
			c := g.Counter(func(ctx Context, specs Specs) (float64, any) { attacked = true; return -1.0, nil })
			c.Attacks(w)
			res, err := g.Evaluate(NewContext())
			if err != nil {
				t.Fatal(err)
			}
			// a = 0.25, unattacked in effect → verdict 2·0.25 − 1 = −0.5, but still active.
			if len(res) != 1 || res[0].Verdict != -0.5 || !attacked {
				t.Fatalf("expected active result -0.5 with attackers evaluated, got %+v (attacked=%v)", res, attacked)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
