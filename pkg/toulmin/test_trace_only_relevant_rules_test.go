//ff:func feature=engine type=engine control=iteration dimension=1
//ff:what TestTraceOnlyRelevantRules — tests that trace contains only rules relevant to each warrant
package toulmin

import (
	"testing"
)

func TestTraceOnlyRelevantRules(t *testing.T) {
	warrantX := func(claim any, ground any, backing Backing) (bool, any) { return true, nil }
	unrelatedDefeater := func(claim any, ground any, backing Backing) (bool, any) { return true, nil }
	g := NewGraph("test")
	wA := g.Warrant(WarrantA, nil, 1.0)
	wX := g.Warrant(warrantX, nil, 1.0)
	ud := g.Defeater(unrelatedDefeater, nil, 1.0)
	rB := g.Rebuttal(RebuttalB, nil, 1.0)
	g.Defeat(rB, wA)
	g.Defeat(ud, wX)
	results, err := g.Evaluate(nil, nil, EvalOption{Trace: true})
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
