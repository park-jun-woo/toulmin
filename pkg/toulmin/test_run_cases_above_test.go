//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesAbove — tests RunCases with VerdictAbove expectation
package toulmin

import "testing"

func TestRunCasesAbove(t *testing.T) {
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 1.0)
	RunCases(t, g, []TestCase{
		{Name: "warrant active", Ground: nil, Expect: VerdictAbove(0)},
	})
}
