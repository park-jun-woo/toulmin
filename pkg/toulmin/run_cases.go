//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what RunCases — runs table-driven test cases against a graph
package toulmin

import "testing"

// RunCases runs each TestCase as a sub-test against the given graph.
func RunCases(t *testing.T, g *Graph, cases []TestCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Helper()
			results, err := g.Evaluate(tc.Context, tc.Option)
			if err != nil {
				t.Fatalf("evaluate error: %v", err)
			}
			if err := tc.Expect(results); err != nil {
				t.Errorf("%v", err)
			}
		})
	}
}
