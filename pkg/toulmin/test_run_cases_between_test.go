//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesBetween — tests RunCases with VerdictBetween expectation
package toulmin

import "testing"

func TestRunCasesBetween(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	d := g.Defeater(DefeaterC, nil, 1.0)
	g.Defeat(r, w)
	g.Defeat(d, r)
	RunCases(t, g, []TestCase{
		{Name: "partial defeat", Ground: nil, Expect: VerdictBetween(0, 0.5)},
	})
}
