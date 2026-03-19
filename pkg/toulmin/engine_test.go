package toulmin

import (
	"fmt"
	"math"
	"testing"
)

func TestWarrantOnly(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0, got %f", results[0].Verdict)
	}
}

func TestWarrantWithDefeater(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 0.0 {
		t.Errorf("expected 0.0, got %f", results[0].Verdict)
	}
}

func TestCompensation(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D1", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D2", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"D1"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	expected := 1.0 / 3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}

func TestStrictWarrant(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Strict,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (strict rejects attack), got %f", results[0].Verdict)
	}
}

func TestCircularAttackError(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "A", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"B"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "B", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"A"},
		Fn:      func(c any, g any, b any) (bool, any) { return true, nil },
	})
	_, err := eng.Evaluate(nil, nil)
	if err == nil {
		t.Fatal("expected error for circular defeat graph")
	}
}

func TestNilFuncGuard(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any, b any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "Ghost", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"W"},
		Fn:      nil,
	})
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (nil attacker ignored), got %f", results[0].Verdict)
	}
}

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


func TestDeepDefeatChainEngine(t *testing.T) {
	fn := func(c any, g any, b any) (bool, any) { return true, nil }
	eng := NewEngine()
	eng.Register(RuleMeta{Name: "W", Qualifier: 1.0, Strength: Defeasible, Fn: fn})
	prev := "W"
	for i := 1; i <= 150; i++ {
		name := fmt.Sprintf("D%d", i)
		eng.Register(RuleMeta{Name: name, Qualifier: 1.0, Strength: Defeater, Defeats: []string{prev}, Fn: fn})
		prev = name
	}
	results, err := eng.Evaluate(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	v := results[0].Verdict
	if math.IsNaN(v) || math.IsInf(v, 0) {
		t.Errorf("verdict should be finite, got %f", v)
	}
}
