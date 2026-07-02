//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunCases — tests RunCases for empty cases (zero iterations) and multiple cases (multiple iterations)
package toulmin

import "testing"

func TestRunCases(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"Empty", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			RunCases(t, g, []TestCase{})
		}},
		{"Multiple", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			RunCases(t, g, []TestCase{
				{Name: "case1", Context: nil, Expect: VerdictAbove(0)},
				{Name: "case2", Context: NewContext(), Expect: VerdictAbove(0)},
			})
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
