//ff:type feature=engine type=model
//ff:what Expectation — verdict assertion function for test cases
package toulmin

// Expectation checks evaluation results against expected conditions.
// Returns nil if the expectation is met, or an error describing the mismatch.
type Expectation func([]EvalResult) error
