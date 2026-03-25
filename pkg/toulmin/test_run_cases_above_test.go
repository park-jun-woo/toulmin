//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesAbove — tests RunCases with VerdictAbove expectation
package toulmin

import "testing"

func TestRunCasesAbove(t *testing.T) {
	g := NewGraph("test")
	g.Rule(WarrantA)
	RunCases(t, g, []TestCase{
		{Name: "warrant active", Context: nil, Expect: VerdictAbove(0)},
	})
}
