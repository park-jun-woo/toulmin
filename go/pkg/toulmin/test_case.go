//ff:type feature=engine type=model
//ff:what TestCase — table-driven test case for policy evaluation
package toulmin

// TestCase defines a single test case for RunCases.
type TestCase struct {
	Name    string      // sub-test name for t.Run
	Context Context     // passed to Evaluate
	Option  EvalOption  // evaluation options (zero value for defaults)
	Expect  Expectation // verdict assertion
}
