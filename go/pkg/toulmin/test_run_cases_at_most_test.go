//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesAtMost — tests RunCases with VerdictAtMost expectation
package toulmin

import "testing"

func TestRunCasesAtMost(t *testing.T) {
	g := NewGraph("test")
	w := g.Rule(WarrantA)
	r := g.Counter(RebuttalB)
	r.Attacks(w)
	RunCases(t, g, []TestCase{
		{Name: "fully rebutted", Context: nil, Expect: VerdictAtMost(0)},
	})
}
