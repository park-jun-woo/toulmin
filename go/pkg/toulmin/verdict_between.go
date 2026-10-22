//ff:func feature=engine type=engine control=sequence
//ff:what VerdictBetween — returns Expectation that checks lo < verdict <= hi
package toulmin

import "fmt"

// VerdictBetween returns an Expectation that passes when lo < verdict <= hi.
func VerdictBetween(lo, hi float64) Expectation {
	return func(results []EvalResult) error {
		if len(results) == 0 {
			return fmt.Errorf("VerdictBetween(%v, %v): expected results, got none", lo, hi)
		}
		v := results[0].Verdict
		if v <= lo || v > hi {
			return fmt.Errorf("VerdictBetween(%v, %v): expected %v < verdict <= %v, got %f", lo, hi, lo, hi, v)
		}
		return nil
	}
}
