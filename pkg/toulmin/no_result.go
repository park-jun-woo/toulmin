//ff:func feature=engine type=engine control=sequence
//ff:what NoResult — Expectation that checks results are empty
package toulmin

import "fmt"

// NoResult is an Expectation that passes when no warrant is active (empty results).
var NoResult Expectation = func(results []EvalResult) error {
	if len(results) != 0 {
		return fmt.Errorf("NoResult: expected no results, got %d", len(results))
	}
	return nil
}
