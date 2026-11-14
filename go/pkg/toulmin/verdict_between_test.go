//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestVerdictBetween — tests VerdictBetween for empty-results, at/below-lo, above-hi, and success branches
package toulmin

import "testing"

func TestVerdictBetween(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"EmptyResults", func(t *testing.T) {
			exp := VerdictBetween(0, 1)
			err := exp(nil)
			if err == nil {
				t.Fatal("expected error for empty results")
			}
		}},
		{"AtOrBelowLo", func(t *testing.T) {
			exp := VerdictBetween(0.2, 0.8)
			err := exp([]EvalResult{{Verdict: 0.2}})
			if err == nil {
				t.Fatal("expected error when verdict <= lo")
			}
		}},
		{"AboveHi", func(t *testing.T) {
			exp := VerdictBetween(0.2, 0.8)
			err := exp([]EvalResult{{Verdict: 0.9}})
			if err == nil {
				t.Fatal("expected error when verdict > hi")
			}
		}},
		{"Success", func(t *testing.T) {
			exp := VerdictBetween(0.2, 0.8)
			err := exp([]EvalResult{{Verdict: 0.5}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
