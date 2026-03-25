//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestTraceOnlyRelevantRules — tests that trace contains only rules relevant to each warrant
package toulmin

import (
	"testing"
)

func TestTraceOnlyRelevantRules(t *testing.T) {
	warrantX := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	unrelatedDefeater := func(ctx Context, specs Specs) (bool, any) { return true, nil }
	g := NewGraph("test")
	wA := g.Rule(WarrantA)
	wX := g.Rule(warrantX)
	ud := g.Except(unrelatedDefeater)
	rB := g.Counter(RebuttalB)
	rB.Attacks(wA)
	ud.Attacks(wX)
	results, err := g.Evaluate(nil, EvalOption{Trace: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	traceA := results[0].Trace
	for _, te := range traceA {
		if te.Name != "WarrantA" && te.Name != "RebuttalB" {
			t.Errorf("WarrantA trace contains unrelated rule: %s", te.Name)
		}
	}
}
