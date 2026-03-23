//ff:func feature=engine type=engine control=sequence
//ff:what TestBackingSameFunc — tests same func with different backings creates separate rules
package toulmin

import (
	"testing"
)

func TestBackingSameFunc(t *testing.T) {
	isInRole := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	g.Warrant(isInRole, &testBacking{Value: "admin"}, 1.0)
	g.Warrant(isInRole, &testBacking{Value: "editor"}, 1.0)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results (same func, different backing), got %d", len(results))
	}
}
