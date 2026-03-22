//ff:func feature=engine type=engine control=sequence
//ff:what TestEngineGraphBuilderConsistency — tests that Engine and Graph builder produce same verdict
package toulmin

import (
	"math"
	"testing"
)

func TestEngineGraphBuilderConsistency(t *testing.T) {
	w := func(c any, g any, b any) (bool, any) { return true, nil }
	r := func(c any, g any, b any) (bool, any) { return true, nil }

	eng := NewEngine()
	eng.Register(RuleMeta{Name: "w", Qualifier: 1.0, Strength: Defeasible, Fn: w})
	eng.Register(RuleMeta{Name: "r", Qualifier: 0.8, Strength: Defeasible, Defeats: []string{"w"}, Fn: r})
	engResults, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("engine error: %v", err)
	}

	g := NewGraph("test")
	wRule := g.Warrant(w, nil, 1.0)
	rRule := g.Rebuttal(r, nil, 0.8)
	g.Defeat(rRule, wRule)
	gbResults, err := g.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("graph builder error: %v", err)
	}

	if len(engResults) != 1 || len(gbResults) != 1 {
		t.Fatalf("expected 1 result each, got %d and %d", len(engResults), len(gbResults))
	}
	if math.Abs(engResults[0].Verdict-gbResults[0].Verdict) > 0.001 {
		t.Errorf("verdict mismatch: engine=%f, graphbuilder=%f", engResults[0].Verdict, gbResults[0].Verdict)
	}
}
