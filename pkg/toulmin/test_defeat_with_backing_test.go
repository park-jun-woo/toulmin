//ff:func feature=engine type=engine control=sequence
//ff:what TestDefeatWithBacking — tests defeat chain with backing parameters
package toulmin

import (
	"math"
	"testing"
)

func TestDefeatWithBacking(t *testing.T) {
	isIPInList := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	isAuth := func(claim any, ground any, backing any) (bool, any) { return true, nil }
	g := NewGraph("test")
	auth := g.Warrant(isAuth, nil, 1.0)
	blocked := g.Rebuttal(isIPInList, "blocklist", 1.0)
	allowed := g.Defeater(isIPInList, "whitelist", 1.0)
	g.Defeat(blocked, auth)
	g.Defeat(allowed, blocked)
	results, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f (whitelist defeats blocklist), got %f", expected, results[0].Verdict)
	}
}
