//ff:func feature=engine type=engine control=sequence
//ff:what VerdictAtMost — returns Expectation that checks verdict <= threshold
package toulmin

import "fmt"

// VerdictAtMost returns an Expectation that passes when the first result's verdict <= v.
func VerdictAtMost(v float64) Expectation {
	return func(results []EvalResult) error {
		if len(results) == 0 {
			return fmt.Errorf("VerdictAtMost(%v): expected results, got none", v)
		}
		if results[0].Verdict > v {
			return fmt.Errorf("VerdictAtMost(%v): expected verdict <= %v, got %f", v, v, results[0].Verdict)
		}
		return nil
	}
}
