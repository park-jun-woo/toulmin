//ff:func feature=engine type=engine control=sequence
//ff:what TestLegacySignature — tests that legacy 2-arg func signature works
package toulmin

import (
	"testing"
)

func TestLegacySignature(t *testing.T) {
	legacyFn := func(claim any, ground any) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Warrant(legacyFn, nil, 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
