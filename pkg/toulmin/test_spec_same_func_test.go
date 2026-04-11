//ff:func feature=engine type=engine control=sequence
//ff:what TestSpecSameFunc — tests same func with different specs creates separate rules
package toulmin

import (
	"testing"
)

func TestSpecSameFunc(t *testing.T) {
	isInRole := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Rule(isInRole).With(&testSpec{Value: "admin"})
	g.Rule(isInRole).With(&testSpec{Value: "editor"})
	results, err := g.Evaluate(NewContext())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (same func, different specs), got %d", len(results))
	}
}
