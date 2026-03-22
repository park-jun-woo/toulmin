//ff:func feature=engine type=engine control=sequence
//ff:what TestBackingNil — tests that nil backing appears as nil in trace
package toulmin

import (
	"testing"
)

func TestBackingNil(t *testing.T) {
	g := NewGraph("test")
	g.Warrant(WarrantA, nil, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Trace[0].Backing != nil {
		t.Errorf("expected backing nil, got %v", results[0].Trace[0].Backing)
	}
}
