//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestGraphEvaluateInternal — evaluate: cycle error, warrant filtering, inactive skip, Applicability, trace supersede error, full-pass fill error
package toulmin

import (
	"strings"
	"testing"
)

func TestGraphEvaluateInternal(t *testing.T) {
	ok := func(ctx Context, specs Specs) (float64, any) { return 0.0, "ev" }
	off := func(ctx Context, specs Specs) (float64, any) { return -1.0, nil }
	boom := func(ctx Context, specs Specs) (float64, any) { panic("boom") }
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"CycleError", func(t *testing.T) {
			g := NewGraph("cyc")
			a := g.Rule(ok)
			b := g.Counter(off)
			a.Attacks(b)
			b.Attacks(a)
			if _, _, err := g.evaluate(NewContext(), EvalOption{}, false); err == nil {
				t.Fatal("expected cycle error")
			}
		}},
		{"WarrantsOnlyAndInactiveSkipped", func(t *testing.T) {
			g := NewGraph("w")
			w := g.Rule(ok)
			g.Rule(func(ctx Context, specs Specs) (bool, any) { return false, nil }) // inactive warrant (bool rule)
			g.Counter(off).Attacks(w)
			res, ec, err := g.evaluate(NewContext(), EvalOption{Trace: true}, false)
			if err != nil || ec == nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res) != 1 || res[0].Verdict != 0.0 || res[0].Applicability != 0.5 || res[0].Evidence != "ev" || len(res[0].Trace) == 0 {
				t.Fatalf("expected one traced warrant (verdict 0, applicability 0.5), got %+v", res)
			}
		}},
		{"RulePanicDuringWarrant", func(t *testing.T) {
			g := NewGraph("p")
			g.Rule(boom)
			if _, _, err := g.evaluate(NewContext(), EvalOption{}, false); err == nil || !strings.Contains(err.Error(), "boom") {
				t.Fatalf("expected panic error, got %v", err)
			}
		}},
		{"SupersederPanicInTrace", func(t *testing.T) {
			// W is evaluated first; its trace pass forces the superseder S, which panics.
			g := NewGraph("s")
			w := g.Rule(ok)
			g.Rule(boom).Supersedes(w)
			if _, _, err := g.evaluate(NewContext(), EvalOption{Trace: true}, false); err == nil || !strings.Contains(err.Error(), "boom") {
				t.Fatalf("expected superseder panic error, got %v", err)
			}
		}},
		{"FullPassFillError", func(t *testing.T) {
			// The counter attacks an inactive warrant, so only the full pass runs it.
			g := NewGraph("f")
			w := g.Rule(off)
			g.Counter(boom).Attacks(w)
			_, ec, err := g.evaluate(NewContext(), EvalOption{}, true)
			if err == nil || ec == nil {
				t.Fatalf("expected fill error with evalContext, got %v / %v", err, ec)
			}
		}},
		{"FullPassSuccess", func(t *testing.T) {
			g := NewGraph("ok")
			w := g.Rule(off)
			g.Counter(ok).Attacks(w)
			res, ec, err := g.evaluate(NewContext(), EvalOption{}, true)
			if err != nil || len(res) != 0 || !ec.ran[ec.edges[w.id][0]] {
				t.Fatalf("expected the full pass to run the counter, got %v / %+v", err, res)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
