//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesBetween — tests RunCases with VerdictBetween expectation
package toulmin

import "testing"

func TestRunCasesBetween(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	d := g.Except(DefeaterC)
	r.Attacks(w)
	// Half-force Except: rebuttal applicability 0.5 → rebuttal verdict 0 → warrant 1/3.
	// (A full-force Except would gate the rebuttal out entirely → warrant 1.0.)
	d.Attacks(r, 0.5)
	RunCases(t, g, []TestCase{
		{Name: "partial defeat", Context: nil, Expect: VerdictBetween(0, 0.5)},
	})
}
