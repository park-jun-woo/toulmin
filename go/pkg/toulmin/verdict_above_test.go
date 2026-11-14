//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestVerdictAbove — tests VerdictAbove for empty-results, verdict-too-low, and success branches
package toulmin

import "testing"

func TestVerdictAbove(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"EmptyResults", func(t *testing.T) {
			exp := VerdictAbove(0)
			err := exp(nil)
			if err == nil {
				t.Fatal("expected error for empty results")
			}
		}},
		{"TooLow", func(t *testing.T) {
			exp := VerdictAbove(0.5)
			err := exp([]EvalResult{{Verdict: 0.5}})
			if err == nil {
				t.Fatal("expected error when verdict <= threshold")
			}
		}},
		{"Success", func(t *testing.T) {
			exp := VerdictAbove(0.5)
			err := exp([]EvalResult{{Verdict: 0.6}})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		}},
	}
	for _, c := range cases {
		t.Run(c.name, c.run)
	}
}
