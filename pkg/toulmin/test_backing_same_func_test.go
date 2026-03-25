//ff:func feature=engine type=engine control=sequence
//ff:what TestBackingSameFunc — tests same func with different backings creates separate rules
package toulmin

import (
	"testing"
)

func TestBackingSameFunc(t *testing.T) {
	isInRole := func(ctx Context, backing Backing) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(isInRole).Backing(&testBacking{Value: "admin"})
	g.Rule(isInRole).Backing(&testBacking{Value: "editor"})
	results, err := g.Evaluate(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (same func, different backing), got %d", len(results))
	}
}
