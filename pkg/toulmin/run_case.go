//ff:func feature=engine type=engine control=sequence
//ff:what runCase — evaluates a single test case and asserts expectation
package toulmin

import "testing"

// runCase evaluates a single TestCase against the graph and checks the expectation.
func runCase(t *testing.T, g *Graph, tc TestCase) {
	t.Helper()
	ctx := tc.Context
	if ctx == nil {
		ctx = NewContext()
	}
	results, err := g.Evaluate(ctx, tc.Option)
	if err != nil {
		t.Fatalf("evaluate error: %v", err)
	}
	if err := tc.Expect(results); err != nil {
		t.Errorf("%v", err)
	}
}
