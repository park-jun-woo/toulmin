//ff:func feature=engine type=engine control=sequence
//ff:what TestEmptyGraph — tests that empty graph returns no results and no error
package toulmin

import (
	"testing"
)

func TestEmptyGraph(t *testing.T) {
	g := NewGraph("test")
	results, err := g.Evaluate(NewContext())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}
