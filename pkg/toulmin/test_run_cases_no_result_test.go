//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesNoResult — tests RunCases with NoResult expectation
package toulmin

import "testing"

func TestRunCasesNoResult(t *testing.T) {
	inactive := func(ctx Context, backing Backing) (bool, any) { return false, nil }
	g := NewGraph("test")
	g.Rule(inactive)
	RunCases(t, g, []TestCase{
		{Name: "inactive warrant", Context: nil, Expect: NoResult},
	})
}
