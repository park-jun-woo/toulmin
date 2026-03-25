//ff:func feature=engine type=engine control=sequence
//ff:what TestLegacySignature — tests that context-based func signature works
package toulmin

import (
	"testing"
)

func TestLegacySignature(t *testing.T) {
	fn := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(fn)
	results, err := g.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}
