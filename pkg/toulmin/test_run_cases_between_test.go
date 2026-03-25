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
	d.Attacks(r)
	RunCases(t, g, []TestCase{
		{Name: "partial defeat", Context: nil, Expect: VerdictBetween(0, 0.5)},
	})
}
