//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestVerdictAtMost — tests VerdictAtMost for empty-results, verdict-too-high, and success branches
package toulmin

import "testing"

func TestVerdictAtMost(t *testing.T) {
	cases := []struct {
		name string
		run  func(t *testing.T)
	}{
		{"EmptyResults", func(t *testing.T) {
			exp := VerdictAtMost(0)
			err := exp(nil)
			if err == nil {
				t.Fatal("expected error for empty results")
			}
		}},
		{"TooHigh", func(t *testing.T) {
			exp := VerdictAtMost(0.5)
			err := exp([]EvalResult{{Verdict: 0.6}})
			if err == nil {
				t.Fatal("expected error when verdict > threshold")
			}
		}},
		{"Success", func(t *testing.T) {
			exp := VerdictAtMost(0.5)
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
