//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestRunCase — tests runCase for nil/non-nil context, evaluate error, and Expect success/error branches
package toulmin

import (
	"errors"
	"testing"
)

func TestRunCase(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"NilContextSuccess", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			runCase(t, g, TestCase{
				Name:    "nil ctx",
				Context: nil,
				Expect:  VerdictAbove(0),
			})
		}},
		{"NonNilContextSuccess", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)
			runCase(t, g, TestCase{
				Name:    "explicit ctx",
				Context: NewContext(),
				Expect:  VerdictAbove(0),
			})
		}},
		{"EvaluateErrorBranch", func(t *testing.T) {
			f1 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			f2 := func(ctx Context, specs Specs) (bool, any) { return true, nil }
			g := NewGraph("cycle")
			a := g.Rule(f1)
			b := g.Counter(f2)
			a.Attacks(b)
			b.Attacks(a)

			ft := isolatedT()
			done := make(chan struct{})
			go func() {
				defer close(done)
				runCase(ft, g, TestCase{Name: "cycle", Expect: VerdictAbove(0)})
			}()
			<-done
			if !ft.Failed() {
				t.Errorf("expected runCase to fail the test on evaluate error")
			}
		}},
		{"ExpectErrorBranch", func(t *testing.T) {
			g := NewGraph("test")
			g.Rule(WarrantA)

			failingExpect := Expectation(func(results []EvalResult) error {
				return errors.New("expectation mismatch")
			})

			ft := isolatedT()
			done := make(chan struct{})
			go func() {
				defer close(done)
				runCase(ft, g, TestCase{Name: "mismatch", Expect: failingExpect})
			}()
			<-done
			if !ft.Failed() {
				t.Errorf("expected runCase to fail the test on Expect mismatch")
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
