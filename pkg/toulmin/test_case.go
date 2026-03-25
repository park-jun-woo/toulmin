//ff:type feature=engine type=model
//ff:what TestCase — table-driven test case for policy evaluation
package toulmin

// TestCase defines a single test case for RunCases.
type TestCase struct {
	Name   string      // sub-test name for t.Run
	Claim  any         // passed to Evaluate as claim
	Ground any         // passed to Evaluate as ground
	Option EvalOption  // evaluation options (zero value for defaults)
	Expect Expectation // verdict assertion
}
