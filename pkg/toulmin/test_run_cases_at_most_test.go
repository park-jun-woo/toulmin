//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesAtMost — tests RunCases with VerdictAtMost expectation
package toulmin

import "testing"

func TestRunCasesAtMost(t *testing.T) {
	g := NewGraph("test")
	w := g.Warrant(WarrantA, nil, 1.0)
	r := g.Rebuttal(RebuttalB, nil, 1.0)
	g.Defeat(r, w)
	RunCases(t, g, []TestCase{
		{Name: "fully rebutted", Ground: nil, Expect: VerdictAtMost(0)},
	})
}
