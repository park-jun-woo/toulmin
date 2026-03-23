//ff:func feature=engine type=engine control=sequence
//ff:what TestBackingInTrace — tests that backing value appears in trace entry
package toulmin

import (
	"testing"
)

func TestBackingInTrace(t *testing.T) {
	isInRole := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Warrant(isInRole, &testBacking{Value: "admin"}, 1.0)
	results, err := g.EvaluateTrace(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, ok := results[0].Trace[0].Backing.(*testBacking)
	if !ok || b.Value != "admin" {
		t.Errorf("expected backing with value 'admin', got %v", results[0].Trace[0].Backing)
	}
}
