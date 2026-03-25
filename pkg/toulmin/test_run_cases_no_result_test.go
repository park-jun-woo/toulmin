//ff:func feature=engine type=engine control=sequence
//ff:what TestRunCasesNoResult — tests RunCases with NoResult expectation
package toulmin

import "testing"

func TestRunCasesNoResult(t *testing.T) {
	inactive := func(claim any, ground any, backing Backing) (bool, any) { return false, nil }
	g := NewGraph("test")
	g.Warrant(inactive, nil, 1.0)
	RunCases(t, g, []TestCase{
		{Name: "inactive warrant", Ground: nil, Expect: NoResult},
	})
}
