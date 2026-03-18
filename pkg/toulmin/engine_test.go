package toulmin

import (
	"math"
	"testing"
)

func TestWarrantOnly(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
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
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
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
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D1", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D2", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"D1"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	// D2 attacks D1 → D1 weakened → W partially restored
	// raw(D2)=1.0, raw(D1)=0.5, raw(W)=0.667 → verdict(W)≈+0.333
	expected := 1.0/3.0
	if math.Abs(results[0].Verdict-expected) > 0.001 {
		t.Errorf("expected ≈%f, got %f", expected, results[0].Verdict)
	}
}

func TestStrictWarrant(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Strict,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "D", Qualifier: 1.0, Strength: Defeater,
		Defeats: []string{"W"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (strict rejects attack), got %f", results[0].Verdict)
	}
}

func TestCircularAttack(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "A", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"B"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "B", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"A"},
		Fn:      func(c any, g any) (bool, any) { return true, nil },
	})
	results := eng.Evaluate(nil, nil)
	// Both A and B are attackers, so neither is a warrant — no results
	if len(results) != 0 {
		t.Fatalf("expected 0 results (both are attackers), got %d", len(results))
	}
}

func TestNilFuncGuard(t *testing.T) {
	eng := NewEngine()
	eng.Register(RuleMeta{
		Name: "W", Qualifier: 1.0, Strength: Defeasible,
		Fn: func(c any, g any) (bool, any) { return true, nil },
	})
	eng.Register(RuleMeta{
		Name: "Ghost", Qualifier: 1.0, Strength: Defeasible,
		Defeats: []string{"W"},
		Fn:      nil,
	})
	results := eng.Evaluate(nil, nil)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Verdict != 1.0 {
		t.Errorf("expected +1.0 (nil attacker ignored), got %f", results[0].Verdict)
	}
}

func TestEngineGraphBuilderConsistency(t *testing.T) {
	w := func(c any, g any) (bool, any) { return true, nil }
	r := func(c any, g any) (bool, any) { return true, nil }

	eng := NewEngine()
	eng.Register(RuleMeta{Name: "w", Qualifier: 1.0, Strength: Defeasible, Fn: w})
	eng.Register(RuleMeta{Name: "r", Qualifier: 0.8, Strength: Defeasible, Defeats: []string{"w"}, Fn: r})
	engResults := eng.Evaluate(nil, nil)

	g := NewGraph("test").
		Warrant(w, 1.0).
		Rebuttal(r, 0.8).
		Defeat(r, w)
	gbResults := g.Evaluate(nil, nil)

	if len(engResults) != 1 || len(gbResults) != 1 {
		t.Fatalf("expected 1 result each, got %d and %d", len(engResults), len(gbResults))
	}
	if math.Abs(engResults[0].Verdict-gbResults[0].Verdict) > 0.001 {
		t.Errorf("verdict mismatch: engine=%f, graphbuilder=%f", engResults[0].Verdict, gbResults[0].Verdict)
	}
}

func TestParseAnnotation(t *testing.T) {
	lines := []string{
		`//tm:backing "Böhm-Jacopini theorem"`,
	}
	meta := ParseAnnotation(lines)
	if meta.Backing != "Böhm-Jacopini theorem" {
		t.Errorf("backing: expected 'Böhm-Jacopini theorem', got '%s'", meta.Backing)
	}
	if meta.Qualifier != 1.0 {
		t.Errorf("qualifier: expected default 1.0, got %f", meta.Qualifier)
	}
}
